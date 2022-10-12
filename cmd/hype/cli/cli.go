package cli

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/markbates/cleo"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type App struct {
	cleo.Cmd
}

func (cmd *App) Main(ctx context.Context, pwd string, args []string) error {
	// if err := cmd.init(pwd); err != nil {
	// 	return err
	// }

	if len(args) == 0 {
		return cleo.ErrNoCommand
	}

	if len(pwd) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return plugins.Wrap(cmd, err)
		}
		pwd = wd
	}

	plugs := cmd.ScopedPlugins()

	if len(plugs) == 0 {
		return plugins.Wrap(cmd, cleo.ErrNoCommands)
	}

	c := plugcmd.FindFromArgs(args, plugs)

	if c == nil {
		return plugins.Wrap(cmd, cleo.ErrNoCommands)
	}

	ctx, cancel := cleo.NewContext(ctx)
	defer cancel()

	return c.Main(ctx, pwd, args[1:])
}

func New() *App {
	app := &App{
		Cmd: cleo.Cmd{
			Name: "hype",
			FS:   os.DirFS("."),
		},
	}

	app.Plugins = append(app.Plugins,
		&Marked{
			Cmd: cleo.Cmd{
				Name:    "marked",
				Aliases: []string{"m", "md"},
			},
		},
		&Marked{
			Cmd: cleo.Cmd{
				Name:    "preview",
				Aliases: []string{"p"},
			},
		},
		&Latex{
			Cmd: cleo.Cmd{
				Name:    "latex",
				Aliases: []string{"l"},
			},
		},
		&VSCode{
			Cmd: cleo.Cmd{
				Name:    "vscode",
				Aliases: []string{"code"},
			},
		},
	)

	return app
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
