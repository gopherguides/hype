package commander

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
	"github.com/mattn/go-shellwords"
	"golang.org/x/net/html"
)

var _ hype.Tag = &Cmd{}
var _ hype.SetSourceable = &Cmd{}

type Cmd struct {
	*hype.Node
	Root string
	Args []string
}

func (c *Cmd) Source() (hype.Source, bool) {
	return hype.SrcAttr(c.Attrs())
}

func (c *Cmd) SetSource(src string) {
	c.Set("src", src)
	// panic(src)
	c.work(c.Root, src)
}

func (c *Cmd) String() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, c.StartTag())
	fmt.Fprint(bb, c.Children.String())
	fmt.Fprint(bb, c.EndTag())
	return bb.String()
}

func (c *Cmd) Validate(checks ...hype.ValidatorFn) error {
	return c.Node.Validate(html.ElementNode, checks...)
}

func NewCmd(node *hype.Node, root string) (*Cmd, error) {
	cmd := &Cmd{
		Node: node,
		Root: root,
	}

	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	ats := cmd.Attrs()
	cmd.work(root, ats["src"])
	return cmd, cmd.Validate()
}

func (cmd *Cmd) work(root string, src string) error {
	ex, err := cmd.Get("exec")
	if err != nil {
		return err
	}

	args, err := shellwords.Parse(ex)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("exec is empty")
	}
	cmd.Args = args

	u, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cache := filepath.Join(u, ".hype", runtime.Version(), "commander")
	os.MkdirAll(cache, 0755)

	runDir := filepath.Join(root, src)
	h, _ := hash(runDir)

	cargs := flect.Underscore(cmd.StartTag())

	cfp := filepath.Join(cache, h, cargs) + ".json"
	os.MkdirAll(filepath.Dir(cfp), 0755)

	ats := cmd.Attrs()

	data := Data{
		// "src": src,
		// "pwd": runDir,
	}
	n := cmd.Args[0]

	if ats.HasKeys("data-go") || n == "go" {
		data["go"] = runtime.Version()
	}

	if !ats.HasKeys("no-cache") {

		if _, err := os.Stat(cfp); err == nil {
			x, err := fromCache(cmd, cfp, data)
			if err != nil {
				return err
			}
			(*cmd) = *x
			return nil
		}
	}

	ctx := context.Background()

	var ag []string
	if len(cmd.Args) > 1 {
		ag = cmd.Args[1:]
	}

	res, err := Run(ctx, runDir, n, ag...)
	if err != nil {
		return err
	}

	data["duration"] = res.Duration.String()

	tag := res.Tag(ats, data)
	cmd.Children = append(cmd.Children, tag)

	f, err := os.Create(cfp)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", cfp, err)
	}
	defer f.Close()

	cf := CacheFile{
		Result: res,
		HTML:   []byte(tag.String()),
	}

	w := io.MultiWriter(f)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	err = enc.Encode(cf)

	if err != nil {
		return fmt.Errorf("could not encode %s: %w", cfp, err)
	}

	return nil
}
