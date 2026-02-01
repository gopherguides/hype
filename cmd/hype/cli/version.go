package cli

import (
	"context"
	"fmt"

	"github.com/markbates/cleo"
)

type Version struct {
	cleo.Cmd

	version string
	commit  string
	date    string
}

func NewVersion(version, commit, date string) *Version {
	return &Version{
		Cmd: cleo.Cmd{
			Name:    "version",
			Aliases: []string{"v"},
			Desc:    "print version information",
		},
		version: version,
		commit:  commit,
		date:    date,
	}
}

func (cmd *Version) Main(ctx context.Context, pwd string, args []string) error {
	if err := (&cmd.Cmd).Init(); err != nil {
		return err
	}

	fmt.Fprintf(cmd.Stdout(), "hype version %s (commit: %s, built: %s)\n", cmd.version, cmd.commit, cmd.date)
	return nil
}
