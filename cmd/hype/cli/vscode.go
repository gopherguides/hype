package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
)

type VSCode struct {
	cleo.Cmd

	Timeout time.Duration
	Host    string

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

	p := hype.NewParser(cmd.FS)

	doc, err := p.ParseExecuteFile(ctx, name)
	if err != nil {
		return err
	}

	body, err := doc.Body()
	if err != nil {
		return err
	}

	toc, err := cmd.toc(p, body)
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

	data := map[string]any{
		"title": doc.Title,
		"body":  body.Children().String(),
		"toc":   toc.String(),
	}

	if err := json.NewEncoder(cmd.Stdout()).Encode(data); err != nil {
		return err
	}

	return nil
}

func (cmd *VSCode) toc(p *hype.Parser, body *hype.Body) (hype.Nodes, error) {
	headings := hype.ByType[*hype.Heading](body.Children())

	bb := &bytes.Buffer{}

	for i, h := range headings {
		t := h.Children().String()

		for i := 1; i < h.Level(); i++ {
			fmt.Fprint(bb, "\t")
		}

		fmt.Fprintf(bb, "1. <a href=\"#heading-%d\">%s</a>\n", i, t)

		link := hype.Text(fmt.Sprintf("<a id=\"heading-%d\"></a>%s", i, t))
		h.Nodes = hype.Nodes{link}
	}

	return p.ParseFragment(bb)
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