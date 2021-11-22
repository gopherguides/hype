package commander

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gopherguides/hype"
	"github.com/mattn/go-shellwords"
	"golang.org/x/net/html"
)

var _ hype.Tag = &Cmd{}
var _ hype.SetSourceable = &Cmd{}

type Cmd struct {
	*hype.Node
	Args []string
	Env  []string
}

func (c *Cmd) Source() (hype.Source, bool) {
	return hype.SrcAttr(c.Attrs())
}

func (c *Cmd) SetSource(src string) {
	c.Set("src", src)
}

func (c *Cmd) Finalize(p *hype.Parser) error {
	return c.work(p.Root, c.Attrs()["src"])
}

func (c *Cmd) StartTag() string {
	return c.Node.StartTag() + `<pre class="code-block"><code class="language-plain" language="plain">`
}

func (c *Cmd) EndTag() string {
	return "</code></pre>" + c.Node.EndTag()
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

func NewCmd(node *hype.Node) (*Cmd, error) {
	cmd := &Cmd{
		Node: node,
	}

	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	ats := cmd.Attrs()

	if env, ok := ats["environ"]; ok {
		cmd.Env = strings.Split(env, ",")
	}

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

	ats := cmd.Attrs()

	data := Data{}

	e := strings.Join(cmd.Env, ",")
	e = strings.TrimSpace(e)
	if len(e) > 0 {
		data["env"] = e
	}

	n := cmd.Args[0]

	if ats.HasKeys("data-go") || n == "go" {
		data["go"] = runtime.Version()
	}

	if !ats.HasKeys("no-cache") {

		if err := cache.Retrieve(cmd, data); err == nil {
			return nil
		}

	}

	ctx := context.Background()

	var ag []string
	if len(cmd.Args) > 1 {
		ag = cmd.Args[1:]
	}

	runDir := filepath.Join(root, src)
	res, err := Run(ctx, runDir, cmd.Env, n, ag...)

	if err != nil {
		return err
	}

	data["duration"] = res.Duration.String()

	s, err := res.Out(ats, data)
	if err != nil {
		return err
	}

	cmd.Children = hype.Tags{hype.QuickText(s)}

	// if res.Err != nil {
	// 	return nil
	// }

	if err := cache.Store(cmd, data, res); err != nil {
		return err
	}

	return nil
}
