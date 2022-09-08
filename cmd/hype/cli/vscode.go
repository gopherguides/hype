package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/atomx"
	"github.com/markbates/cleo"
)

type VSCode struct {
	cleo.Cmd

	Timeout time.Duration
	Host    string
	Parser  *hype.Parser // If nil, a default parser is used.

	flags *flag.FlagSet
}

func (cmd *VSCode) Flags() (*flag.FlagSet, error) {
	if err := cmd.validate(); err != nil {
		return nil, err
	}

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("vscode", flag.ContinueOnError)
	cmd.flags.SetOutput(cmd.Stderr())
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", DefaultTimeout(), "timeout for execution")
	cmd.flags.StringVar(&cmd.Host, "host", "", "host to serve on")

	return cmd.flags, nil
}

func (cmd *VSCode) Main(ctx context.Context, pwd string, args []string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	flags, err := cmd.Flags()
	if err != nil {
		return err
	}

	if err := flags.Parse(args); err != nil {
		return err
	}

	args = flags.Args()

	if len(args) < 1 {
		return fmt.Errorf("no filename specified")
	}

	if len(cmd.Host) == 0 {
		return fmt.Errorf("no host specified")
	}

	path := args[0]

	pwd = filepath.Dir(path)
	name := filepath.Base(path)

	cmd.FS = os.DirFS(pwd)

	err = WithTimeout(ctx, cmd.Timeout, func(ctx context.Context) error {
		return WithinDir(pwd, func() error {
			return cmd.execute(ctx, pwd, name)
		})
	})

	if err != nil {
		return err
	}
	return nil
}

func (cmd *VSCode) execute(ctx context.Context, pwd string, name string) error {
	err := cmd.validate()
	if err != nil {
		return err
	}

	p := cmd.Parser
	if p == nil {
		p = hype.NewParser(cmd.FS)
	}

	ncn := p.NodeParsers[atomx.Code]
	if ncn == nil {
		ncn = hype.NewCodeNodes
	}

	p.NodeParsers[atomx.Code] = func(p *hype.Parser, el *hype.Element) (hype.Nodes, error) {
		nodes, err := ncn(p, el)
		if err != nil {
			return nil, err
		}

		src, ok := el.Get("src")
		if !ok {
			return nodes, nil
		}

		vel := hype.NewEl("vscode", nil)
		vs := &vscode{
			Element: vel,
		}

		vs.Set("src", src)

		nodes = append(nodes, vs)

		return nodes, nil

	}

	doc, err := p.ParseExecuteFile(ctx, name)
	if err != nil {
		return err
	}

	body, err := doc.Body()
	if err != nil {
		return err
	}

	images := hype.ByType[*hype.Image](body.Children())
	for _, i := range images {
		src, _ := i.Get("src")
		src = filepath.Join(cmd.Host, src)
		if err := i.Set("src", src); err != nil {
			return err
		}
	}

	tc, err := cmd.toc(p, body)
	if err != nil {
		return err
	}

	data := map[string]any{
		"title": doc.Title,
		"body":  body.Children().String(),
		"toc":   tc,
	}

	if err := json.NewEncoder(cmd.Stdout()).Encode(data); err != nil {
		return err
	}

	return nil
}

func (cmd *VSCode) toc(p *hype.Parser, body *hype.Body) (string, error) {
	toc, err := hype.GenerateToC(p, body.Children())
	if err != nil {
		return "", err
	}

	headings := hype.ByType[*hype.Heading](body.Children())

	for i, h := range headings {
		x := h.Children().String()
		link := hype.Text(fmt.Sprintf("<a id=\"heading-%d\"></a>%s", i, x))
		h.Nodes = hype.Nodes{link}
	}

	tc := fmt.Sprintf("<div id=\"menu\">\n%s\n</div>\n", toc.String())

	return tc, nil
}

func (cmd *VSCode) validate() error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.FS == nil {
		cmd.FS = os.DirFS(".")
	}

	if cmd.Timeout == 0 {
		cmd.Timeout = DefaultTimeout()
	}

	return nil
}

type vscode struct {
	*hype.Element
}

func (v *vscode) PostExecute(ctx context.Context, doc *hype.Document, err error) error {
	if err != nil {
		return nil
	}

	src, _ := v.Get("src")

	src = strings.Split(src, "#")[0]
	v.Nodes = append(v.Nodes, hype.Text(fmt.Sprintf("<a href=\"%s\" class=\"hype-file-open\"><i class=\"bi bi-file-earmark-code-fill\"></i> %s</a>", src, src)))
	return nil
}
