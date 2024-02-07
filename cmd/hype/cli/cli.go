package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type App struct {
	cleo.Cmd

	Parser *hype.Parser

	once    sync.Once
	initErr error
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

	if err := cmd.init(pwd, args); err != nil {
		return err
	}

	return cmd.execute(ctx, pwd, args)
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

func (cmd *App) execute(ctx context.Context, pwd string, args []string) error {
	if cmd == nil {
		return plugins.Wrap(cmd, fmt.Errorf("app is nil"))
	}

	plugs := plugins.Plugins{}
	for _, c := range cmd.SubCommands() {
		plugs = append(plugs, c)
	}

	c := plugcmd.FindFromArgs(args, plugs)

	if c == nil {
		return plugins.Wrap(cmd, cleo.ErrNoCommands)
	}

	err := c.Main(ctx, pwd, args[1:])
	if err == nil {
		return nil
	}

	return err
}

func (cmd *App) init(pwd string, args []string) error {
	if cmd == nil {
		return plugins.Wrap(cmd, fmt.Errorf("app is nil"))
	}

	cmd.once.Do(func() {
		if err := cleo.Init(&cmd.Cmd, pwd); err != nil {
			cmd.initErr = plugins.Wrap(cmd, err)
		}

		for _, c := range cmd.SubCommands() {
			if pc, ok := c.(ParserCommander); ok {
				if err := pc.SetParser(cmd.Parser); err != nil {
					cmd.initErr = plugins.Wrap(cmd, err)
					break
				}
			}
		}
	})

	return cmd.initErr
}

func New(root string) *App {
	cab := os.DirFS(root)

	p := hype.NewParser(cab)

	m := &Marked{
		Cmd: cleo.Cmd{
			Name:    "marked",
			Aliases: []string{"m", "md"},
			Desc:    "outputs for the Marked app (https://marked2app.com/)",
		},
		Parser: p,
	}

	e := &Export{
		Cmd: cleo.Cmd{
			Name:    "export",
			Aliases: []string{"export", "e"},
			Desc:    "export the document to a different format (markdown,json,html,etc...)",
		},
		Parser: p,
	}

	mp := &Marked{
		Cmd: cleo.Cmd{
			Name:    "preview",
			Aliases: []string{"p"},
			Desc:    "outputs HTML for previwing document in a browser",
		},
		Parser: p,
	}

	sl := &Slides{
		Cmd: cleo.Cmd{
			Name:    "slides",
			Aliases: []string{"s"},
			Desc:    "outputs slide format for presentation",
		},
		Parser: p,
	}

	app := &App{
		Cmd: cleo.Cmd{
			Name: "hype",
			FS:   cab,
			Commands: map[string]cleo.Commander{
				"marked":  m,
				"preview": mp,
				"slides":  sl,
				"export":  e,
			},
		},
		Parser: p,
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
