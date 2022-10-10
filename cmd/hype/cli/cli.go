package cli

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/markbates/cleo"
)

type App = cleo.Cmd

func New() *App {
	app := &App{
		FS: os.DirFS("."),
	}

	app.Plugins = append(app.Plugins,
		&Marked{
			Cmd: cleo.Cmd{
				Name: "marked",
			},
		},
		&Marked{
			Cmd: cleo.Cmd{
				Name: "preview",
			},
		},
		&Latex{
			Cmd: cleo.Cmd{
				Name: "latex",
			},
		},
		&VSCode{
			Cmd: cleo.Cmd{
				Name: "vscode",
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
