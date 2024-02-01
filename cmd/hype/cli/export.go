package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
	"github.com/markbates/plugins"
)

type Export struct {
	cleo.Cmd

	// a folder containing all chapters of a book, for example
	ContextPath string
	File        string        // optional file name to preview
	Timeout     time.Duration // default: 5s
	Parser      *hype.Parser  // If nil, a default parser is used.
	Section     int           // default: 1
	Verbose     bool          // default: false
	Format      string        // default:markdown

	flags *flag.FlagSet
}

func (cmd *Export) WithPlugins(fn plugins.Feeder) {
	if cmd == nil {
		return
	}

	cmd.Lock()
	defer cmd.Unlock()
	cmd.Feeder = fn
}

func (cmd *Export) ScopedPlugins() plugins.Plugins {
	if cmd == nil {
		return nil
	}

	type marker interface {
		MarkedPlugin()
	}

	plugs := cmd.Cmd.ScopedPlugins()

	res := make(plugins.Plugins, 0, len(plugs))
	for _, p := range plugs {
		if _, ok := p.(marker); ok {
			res = append(res, p)
		}
	}

	return res
}

func (cmd *Export) SetParser(p *hype.Parser) error {
	if cmd == nil {
		return fmt.Errorf("marked is nil")
	}

	cmd.Lock()
	defer cmd.Unlock()

	cmd.Parser = p
	return nil
}

func (cmd *Export) Flags() (*flag.FlagSet, error) {
	if err := cmd.validate(); err != nil {
		return nil, err
	}

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("marked", flag.ContinueOnError)
	cmd.flags.SetOutput(io.Discard)
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", DefaultTimeout(), "timeout for execution")
	cmd.flags.StringVar(&cmd.ContextPath, "context", cmd.ContextPath, "a folder containing all chapters of a book, for example")
	cmd.flags.StringVar(&cmd.File, "f", "module.md", "optional file name to preview, if not provided, defaults to module.md")
	cmd.flags.IntVar(&cmd.Section, "section", 0, "")
	cmd.flags.BoolVar(&cmd.Verbose, "v", false, "enable verbose output for debugging")
	cmd.flags.StringVar(&cmd.Format, "format", "markdown", "content type to export to (markdown, html, body, etc...)  See documentation for more options")

	return cmd.flags, nil
}

func (cmd *Export) Main(ctx context.Context, pwd string, args []string) error {
	err := cmd.main(ctx, pwd, args)
	if err == nil {
		return nil
	}

	cmd.Lock()
	to := cmd.Timeout
	if to == 0 {
		to = DefaultTimeout()
		cmd.Timeout = to
	}
	cmd.Unlock()

	ctx, cancel := cleo.ContextWithTimeout(ctx, to)
	defer cancel()

	var mu sync.Mutex

	go func() {
		mu.Lock()
		err = plugins.Wrap(cmd, err)
		mu.Unlock()
		cancel()
	}()

	<-ctx.Done()

	mu.Lock()
	defer mu.Unlock()
	return err
}

func (cmd *Export) main(ctx context.Context, pwd string, args []string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	mp := os.Getenv("MARKED_PATH")

	pwd = filepath.Dir(mp)

	if err := cleo.Init(&cmd.Cmd, pwd); err != nil {
		return err
	}

	flags, err := cmd.Flags()
	if err != nil {
		return err
	}

	if err := flags.Parse(args); err != nil {
		return err
	}
	if cmd.Verbose {
		// enable debugging
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	err = WithTimeout(ctx, cmd.Timeout, func(ctx context.Context) error {
		// TODO Document what this does
		if mo, ok := os.LookupEnv("MARKED_ORIGIN"); ok {
			pwd = mo
		}

		return WithinDir(pwd, func() error {
			return cmd.execute(ctx, pwd)
		})
	})

	if err != nil {
		return err
	}

	return nil
}

func (cmd *Export) execute(ctx context.Context, pwd string) error {
	err := cmd.validate()
	if err != nil {
		return err
	}

	if cmd.FS == nil {
		cmd.FS = os.DirFS(pwd)
	}

	// TODO Document what this does
	mp := os.Getenv("MARKED_PATH")

	slog.Debug("execute", "pwd", pwd, "file", cmd.File, "MARKED_PATH", mp)
	p := cmd.Parser

	if p == nil {
		p = hype.NewParser(cmd.FS)
	}

	if p.Section == 0 {
		p.Section = 1
	}

	if cmd.Section > 0 {
		p.Section = cmd.Section
	}

	p.Root = filepath.Dir(mp)

	if len(cmd.File) > 0 {
		f, err := cmd.FS.Open(cmd.File)
		if err != nil {
			return err
		}
		defer f.Close()

		cmd.IO.In = f
	}

	doc, err := p.Parse(cmd.Stdin())
	if err != nil {
		return err
	}

	switch cmd.Format {
	case "markdown":
		fmt.Fprintln(cmd.Stdout(), doc.MD())
	case "html":
		fmt.Fprintln(cmd.Stdout(), doc.String())

	// TODO: Implement this
	case "body":
		//fmt.Fprintln(cmd.Stdout(), doc.Body())
		return fmt.Errorf("body format not implemented")
	}

	return nil

}

func (cmd *Export) validate() error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.Timeout == 0 {
		cmd.Timeout = DefaultTimeout()
	}

	return nil
}
