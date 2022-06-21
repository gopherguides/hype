package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
)

type Marked struct {
	cleo.Cmd

	Timeout time.Duration // default: 5s

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
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", cmd.DefaultTimeout(), "timeout for execution")

	return cmd.flags, nil
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
		cmd.Timeout = cmd.DefaultTimeout()
	}

	return nil
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

	cltx, cancel := cleo.ContextWithTimeout(ctx, cmd.Timeout)
	defer cancel()

	go func() {
		defer cancel()

		err := cmd.execute(cltx, pwd)
		if err != nil {
			cltx.SetErr(err)
		}

	}()

	<-cltx.Done()

	err = cltx.Err()
	if !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (cmd *Marked) execute(ctx context.Context, pwd string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	if mo, ok := os.LookupEnv("MARKED_ORIGIN"); ok {
		pwd = mo
	}

	if len(pwd) > 0 {
		opwd, err := os.Getwd()
		if err != nil {
			return err
		}
		defer os.Chdir(opwd)

		if err := os.Chdir(pwd); err != nil {
			return err
		}
	}

	p := hype.NewParser(cmd.FS)

	p.Section = 1

	if mp, ok := os.LookupEnv("MARKED_PATH"); ok {
		if sec, err := SectionFromPath(mp); err == nil {
			p.Section = sec
		}
	}

	p.PreParsers = append(p.PreParsers, &Binding{
		Binder: flect.New("book"),
		Ident:  flect.New("chapter"),
	})

	doc, err := p.ParseExecute(ctx, cmd.Stdin())
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.Stdout(), doc.String())

	return nil

}

func (cmd *Marked) DefaultTimeout() time.Duration {
	return time.Second * 5
}
