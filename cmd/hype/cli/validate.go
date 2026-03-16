package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
)

type Validate struct {
	cleo.Cmd

	File    string
	Timeout time.Duration
	Parser  *hype.Parser
	Verbose bool
	Exec    bool
	Format  string

	flags *flag.FlagSet
	mu    sync.RWMutex
}

func (cmd *Validate) SetParser(p *hype.Parser) error {
	if cmd == nil {
		return fmt.Errorf("validate is nil")
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	cmd.Parser = p
	return nil
}

func (cmd *Validate) Flags(stderr io.Writer) (*flag.FlagSet, error) {
	usage := `
Usage: hype validate [options]

Examples:
	hype validate -f document.md
	hype validate -f document.md --exec
	hype validate -f document.md -v
	hype validate -f document.md --format=json
`

	if err := cmd.validate(); err != nil {
		return nil, err
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("validate", flag.ContinueOnError)
	cmd.flags.SetOutput(stderr)
	cmd.flags.StringVar(&cmd.File, "f", "hype.md", "file to validate")
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", DefaultTimeout, "timeout for execution, defaults to 30 seconds (30s)")
	cmd.flags.BoolVar(&cmd.Verbose, "v", false, "enable verbose output")
	cmd.flags.BoolVar(&cmd.Exec, "exec", false, "also validate code execution")
	cmd.flags.StringVar(&cmd.Format, "format", "text", "output format: text, json")

	cmd.flags.Usage = func() {
		fmt.Fprintf(stderr, "Usage of %s:\n", os.Args[0])
		cmd.flags.PrintDefaults()
		fmt.Fprintln(stderr, usage)
	}

	return cmd.flags, nil
}

func (cmd *Validate) Main(ctx context.Context, pwd string, args []string) error {
	cmd.mu.Lock()
	to := cmd.Timeout
	if to == 0 {
		to = DefaultTimeout
		cmd.Timeout = to
	}
	cmd.mu.Unlock()

	return cmd.main(ctx, pwd, args)
}

func (cmd *Validate) main(ctx context.Context, pwd string, args []string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	if err := (&cmd.Cmd).Init(); err != nil {
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
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	return WithTimeout(ctx, cmd.Timeout, func(ctx context.Context) error {
		return WithinDir(pwd, func() error {
			return cmd.execute(ctx, pwd)
		})
	})
}

func (cmd *Validate) execute(ctx context.Context, pwd string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	if cmd.FS == nil {
		cmd.FS = os.DirFS(pwd)
	}

	slog.Debug("validate", "pwd", pwd, "file", cmd.File)
	fileDir := filepath.Dir(cmd.File)
	fileName := filepath.Base(cmd.File)

	parserFS := cmd.FS
	if fileDir != "." && fileDir != "" {
		subFS, err := fs.Sub(cmd.FS, fileDir)
		if err != nil {
			return fmt.Errorf("failed to create sub filesystem for %s: %w", fileDir, err)
		}
		parserFS = subFS
	}

	p := cmd.Parser
	if p == nil {
		p = hype.NewParser(parserFS)
	} else {
		p.FS = parserFS
	}

	p.Root = filepath.Join(pwd, fileDir)

	doc, err := p.ParseFile(fileName)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	result := hype.Validate(ctx, doc, hype.ValidateOptions{Exec: cmd.Exec})

	out := cmd.Stdout()

	switch cmd.Format {
	case "json":
		b, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(b))
	default:
		for _, issue := range result.Issues {
			fmt.Fprintln(out, issue.String())
		}
		if len(result.Issues) > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintln(out, result.Summary())
	}

	if result.HasErrors() {
		return fmt.Errorf("validation failed: %s", result.Summary())
	}

	return nil
}

func (cmd *Validate) validate() error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	if cmd.Timeout == 0 {
		cmd.Timeout = DefaultTimeout
	}

	return nil
}
