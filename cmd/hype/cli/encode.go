package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
)

type Encode struct {
	cleo.Cmd

	File      string        // optional file name to preview
	Timeout   time.Duration // default: 5s
	Parser    *hype.Parser  // If nil, a default parser is used.
	ParseOnly bool          // if true, only parse the file and exit

	flags *flag.FlagSet

	mu sync.RWMutex
}

func (cmd *Encode) SetParser(p *hype.Parser) error {
	if cmd == nil {
		return fmt.Errorf("encode is nil")
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	cmd.Parser = p
	return nil
}

func (cmd *Encode) Flags() (*flag.FlagSet, error) {
	if cmd == nil {
		return nil, fmt.Errorf("marked is nil")
	}

	cmd.mu.Lock()
	defer cmd.mu.Unlock()

	if cmd.flags != nil {
		return cmd.flags, nil
	}

	cmd.flags = flag.NewFlagSet("marked", flag.ContinueOnError)
	cmd.flags.SetOutput(io.Discard)
	cmd.flags.BoolVar(&cmd.ParseOnly, "p", cmd.ParseOnly, "if true, only parse the file and exit")
	cmd.flags.DurationVar(&cmd.Timeout, "timeout", DefaultTimeout(), "timeout for execution")
	cmd.flags.StringVar(&cmd.File, "f", cmd.File, "optional file name to preview")

	return cmd.flags, nil
}

func (cmd *Encode) Main(ctx context.Context, pwd string, args []string) error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	cmd.mu.Lock()
	if cmd.FS == nil {
		cmd.FS = os.DirFS(pwd)
	}

	if cmd.Parser == nil {
		cmd.Parser = hype.NewParser(cmd.FS)
	}

	p := cmd.Parser

	cmd.mu.Unlock()

	flags, err := cmd.Flags()
	if err != nil {
		return err
	}

	err = flags.Parse(args)
	if err != nil {
		return err
	}

	args = flags.Args()

	if len(args) > 0 {
		fn := args[0]

		doc, err := p.ParseFile(fn)
		if err != nil {
			return err
		}

		if !cmd.ParseOnly {
			err = doc.Execute(ctx)
			if err != nil {
				return err
			}
		}

		enc := json.NewEncoder(cmd.Stdout())
		enc.SetIndent("", "  ")
		return enc.Encode(doc)
	}

	return nil
}
