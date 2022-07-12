package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/markbates/cleo"
)

type Marked struct {
	cleo.Cmd

	// a folder containing all chapters of a book, for example
	ContextPath string
	File        string        // optional file name to preview
	Timeout     time.Duration // default: 5s

	flags *flag.FlagSet
}

func (cmd *Marked) Flags() (*flag.FlagSet, error) {
	if err := cmd.validate(); err != nil {
		return nil, err
	}

	cmd.Lock()
	defer cmd.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("marked", flag.ContinueOnError)
	cmd.flags.SetOutput(cmd.Stderr())
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", DefaultTimeout(), "timeout for execution")
	cmd.flags.StringVar(&cmd.ContextPath, "context", cmd.ContextPath, "a folder containing all chapters of a book, for example")
	cmd.flags.StringVar(&cmd.File, "f", cmd.File, "optional file name to preview")

	return cmd.flags, nil
}

func (cmd *Marked) Main(ctx context.Context, pwd string, args []string) error {
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

	err = WithTimeout(ctx, cmd.Timeout, func(ctx context.Context) error {
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

func (cmd *Marked) execute(ctx context.Context, pwd string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	if cmd.FS == nil {
		cmd.FS = os.DirFS(pwd)
	}

	mp := os.Getenv("MARKED_PATH")

	p, err := NewParser(cmd.FS, cmd.ContextPath, mp)
	if err != nil {
		return err
	}

	if len(cmd.File) > 0 {
		f, err := cmd.FS.Open(cmd.File)
		if err != nil {
			return err
		}
		defer f.Close()

		cmd.IO.In = f
	}

	doc, err := p.ParseExecute(ctx, cmd.Stdin())
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.Stdout(), doc.String())

	return nil

}

func (cmd *Marked) validate() error {
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
