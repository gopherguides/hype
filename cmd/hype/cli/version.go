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

	fmt.Fprintf(cmd.Stdout(), "hype version %s (commit: %s, built: %s)\n", cmd.Info.Version, cmd.Info.Commit, cmd.Info.Date)
	return nil
}
