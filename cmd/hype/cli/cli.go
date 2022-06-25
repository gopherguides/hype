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

	app.Add("preview", &Marked{})
	app.Add("marked", &Marked{})
	app.Add("latex", &Latex{})
	return app
}

func DefaultTimeout() time.Duration {
	return time.Second * 5
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
