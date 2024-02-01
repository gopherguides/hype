package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
)

type TOC struct {
	cleo.Cmd
}

func (cmd *TOC) Main(ctx context.Context, pwd string, args []string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	path := pwd

	if len(args) > 0 {
		path = args[0]
	}

	if cmd.FS == nil {
		cmd.FS = os.DirFS(path)
	}

	fmt.Printf("TODO >> toc.go:32 path %[1]T %[1]v\n", path)

	p := hype.NewParser(cmd.FS)

	docs, err := p.ParseFolder(path)
	if err != nil {
		return fmt.Errorf("error parsing folder, %q: %w", path, err)
	}

	for _, doc := range docs {
		fmt.Fprintln(cmd.Stdout(), doc.Title)
	}

	return nil
}

func (cmd *TOC) validate() error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	// cmd.Lock()
	// defer cmd.Unlock()

	// if cmd.Timeout == 0 {
	// 	cmd.Timeout = DefaultTimeout()
	// }

	return nil
}
