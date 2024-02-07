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
	File    string        // optional file name to preview
	Timeout time.Duration // default: 30s
	Parser  *hype.Parser  // If nil, a default parser is used.
	Verbose bool          // default: false
	Format  string        // default:markdown

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

func (cmd *Export) Flags(stderr io.Writer) (*flag.FlagSet, error) {
	usage := `
Usage: hype export [options]

Examples:
	hype export -format html
	hype export -f README.md -format html
	hype export -f README.md -format markdown -timeout=10s
`

	if err := cmd.validate(); err != nil {
		return nil, err
	}

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("export", flag.ContinueOnError)
	cmd.flags.SetOutput(stderr)
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", DefaultTimeout(), "timeout for execution, defaults to 30 seconds (30s)")
	cmd.flags.StringVar(&cmd.File, "f", "module.md", "optional file name to preview, if not provided, defaults to module.md")
	cmd.flags.BoolVar(&cmd.Verbose, "v", false, "enable verbose output for debugging")
	cmd.flags.StringVar(&cmd.Format, "format", "markdown", "content type to export to: markdown, html")

	cmd.flags.Usage = func() {
		fmt.Fprintf(stderr, "Usage of %s:\n", os.Args[0])
		cmd.flags.PrintDefaults()
		fmt.Fprintln(stderr, usage)
	}

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

	flags, err := cmd.Flags(cmd.Stderr())
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

	p.Root = filepath.Dir(mp)

	doc, err := p.ParseFile(cmd.File)
	if err != nil {
		return err
	}

	if err := doc.Execute(ctx); err != nil {
		return err
	}

	switch cmd.Format {
	case "markdown":
		fmt.Fprintln(cmd.Stdout(), doc.MD())
	case "html":
		fmt.Fprintln(cmd.Stdout(), doc.String())
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", cmd.Format)
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
