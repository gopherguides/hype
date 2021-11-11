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

type Cmd struct {
	*hype.Node
	Root string
	Args []string
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

	ex, err := cmd.Get("exec")
	if err != nil {
		return nil, err
	}

	args, err := shellwords.Parse(ex)
	if err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("exec is empty")
	}
	cmd.Args = args

	u, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cache := filepath.Join(u, ".hype", runtime.Version(), "commander")
	os.MkdirAll(cache, 0755)

	dir, _ := cmd.Get("dir")
	runDir := filepath.Join(root, dir)
	h, _ := hash(runDir)

	cargs := flect.Underscore(cmd.StartTag())

	cfp := filepath.Join(cache, h, cargs) + ".json"
	os.MkdirAll(filepath.Dir(cfp), 0755)

	ats := cmd.Attrs()

	data := Data{}
	n := cmd.Args[0]

	if ats.HasKeys("data-go") || n == "go" {
		data["go"] = runtime.Version()
	}

	if !ats.HasKeys("no-cache") {
		if _, err := os.Stat(cfp); err == nil {
			return fromCache(cmd, cfp, data)
		}
	}

	ctx := context.Background()

	var ag []string
	if len(cmd.Args) > 1 {
		ag = cmd.Args[1:]
	}

	res, err := Run(ctx, runDir, n, ag...)
	if err != nil {
		return nil, err
	}
	data["duration"] = res.Duration.String()

	tag := res.Tag(ats, data)
	cmd.Children = append(cmd.Children, tag)

	f, err := os.Create(cfp)
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", cfp, err)
	}
	defer f.Close()

	cf := CacheFile{
		Result: res,
		HTML:   []byte(tag.String()),
	}

	w := io.MultiWriter(f)
	// w := f

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	err = enc.Encode(cf)

	if err != nil {
		return nil, fmt.Errorf("could not encode %s: %w", cfp, err)
	}

	return cmd, cmd.Validate()
}
