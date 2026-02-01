package cli

import (
	"context"
	"fmt"

	"github.com/markbates/cleo"
)

type Version struct {
	cleo.Cmd

	Info VersionInfo
}

func (cmd *Version) Main(ctx context.Context, pwd string, args []string) error {
	if err := (&cmd.Cmd).Init(); err != nil {
		return err
	}

	fmt.Fprintln(cmd.Stdout(), cmd.Info)
	return nil
}
