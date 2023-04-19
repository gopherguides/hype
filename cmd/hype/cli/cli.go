package cli

import (
	"context"
	"errors"
	"flag"
	"os"
	"time"

	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type App struct {
	cleo.Cmd

	Parser *hype.Parser
}

func (cmd *App) Main(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return cleo.ErrNoCommand
	}

	if len(args) == 1 {
		switch args[0] {
		case "help", "-h", "--help":
			return flag.ErrHelp
		}
	}

	if err := cleo.Init(&cmd.Cmd, pwd); err != nil {
		return err
	}

	c := plugcmd.FindFromArgs(args, cmd.ScopedPlugins())

	if c == nil {
		return plugins.Wrap(cmd, cleo.ErrNoCommands)
	}

	err := c.Main(ctx, pwd, args[1:])
	if err == nil {
		return nil
	}

	// if errors.Is(err, flag.ErrHelp) {
	// 	return nil
	// }

	// if errors.Is(err, cleo.ErrNoCommand) {
	// 	return nil
	// }

	return err
}

func (cmd *App) ScopedPlugins() plugins.Plugins {
	if cmd == nil {
		return nil
	}

	plugs := cmd.Cmd.ScopedPlugins()

	res := make(plugins.Plugins, 0, len(plugs))
	for _, p := range plugs {
		if p != cmd {
			res = append(res, p)
		}
	}

	return res
}

func New(root string) *App {
	// panic(root)
	app := &App{
		Cmd: cleo.Cmd{
			Name: "hype",
			FS:   os.DirFS(root),
		},
	}

	app.Feeder = func() plugins.Plugins {
		return plugins.Plugins{
			&Marked{
				Cmd: cleo.Cmd{
					Name:    "marked",
					Aliases: []string{"m", "md"},
				},
				Parser: app.Parser,
			},
			&Marked{
				Cmd: cleo.Cmd{
					Name:    "preview",
					Aliases: []string{"p"},
				},
				Parser: app.Parser,
			},
			// &Latex{
			// 	Cmd: cleo.Cmd{
			// 		Name:    "latex",
			// 		Aliases: []string{"l"},
			// 	},
			// },
			// &VSCode{
			// 	Cmd: cleo.Cmd{
			// 		Name:    "vscode",
			// 		Aliases: []string{"code"},
			// 	},
			// },
		}
	}

	return app
}

func Garlic(root string) (*App, error) {
	return New(root), nil
}

func DefaultTimeout() time.Duration {
	return time.Second * 30
}

func WithinDir(dir string, f func() error) error {
	if len(dir) == 0 || dir == "." {
		return f()
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(pwd)

	if err := os.Chdir(dir); err != nil {
		return err
	}

	return f()
}

func WithTimeout(ctx context.Context, timeout time.Duration, f func(context.Context) error) error {
	if timeout == 0 {
		timeout = DefaultTimeout()
	}

	cltx, cancel := cleo.ContextWithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		defer cancel()
		err := f(cltx)
		if err != nil {
			cltx.SetErr(err)
		}
	}()

	<-cltx.Done()

	err := cltx.Err()
	if !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
