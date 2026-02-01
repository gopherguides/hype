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

const (
	DefaultTimeout = time.Second * 30
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
		return fmt.Errorf("app is nil")
	}

	plugs := plugins.Plugins{}
	for _, c := range cmd.SubCommands() {
		plugs = append(plugs, c)
	}

	c := plugcmd.FindFromArgs(args, plugs)

	if c == nil {
		return cleo.ErrNoCommands
	}

	err := c.Main(ctx, pwd, args[1:])
	if err == nil {
		return nil
	}

	return err
}

func (cmd *App) init(pwd string, args []string) error {
	if cmd == nil {
		return fmt.Errorf("app is nil")
	}

	cmd.once.Do(func() {
		if err := (&cmd.Cmd).Init(); err != nil {
			cmd.initErr = err
		}

		for _, c := range cmd.SubCommands() {
			if pc, ok := c.(ParserCommander); ok {
				if err := pc.SetParser(cmd.Parser); err != nil {
					cmd.initErr = err
					break
				}
			}
		}
	})

	return cmd.initErr
}

type VersionInfo struct {
	Version string
	Commit  string
	Date    string
}

func New(root string, vi ...VersionInfo) *App {
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

	var versionInfo VersionInfo
	if len(vi) > 0 {
		versionInfo = vi[0]
	} else {
		versionInfo = VersionInfo{Version: "dev", Commit: "none", Date: "unknown"}
	}

	pv := &Preview{
		Cmd: cleo.Cmd{
			Name:    "preview",
			Aliases: []string{"p"},
			Desc:    "live preview server with file watching and auto-reload",
		},
		Parser:  p,
		Version: versionInfo.Version,
	}

	sl := &Slides{
		Cmd: cleo.Cmd{
			Name:    "slides",
			Aliases: []string{"s"},
			Desc:    "outputs slide format for presentation",
		},
		Parser: p,
	}

	bl := &Blog{
		Cmd: cleo.Cmd{
			Name:    "blog",
			Aliases: []string{"b"},
			Desc:    "static blog generator commands (init, build, new)",
		},
	}

	ver := NewVersion(versionInfo.Version, versionInfo.Commit, versionInfo.Date)

	app := &App{
		Cmd: cleo.Cmd{
			Name: "hype",
			FS:   cab,
			Commands: map[string]cleo.Commander{
				"marked":  m,
				"preview": pv,
				"slides":  sl,
				"export":  e,
				"blog":    bl,
				"version": ver,
			},
		},
		Parser: p,
	}

	return app
}

func Garlic(root string) (*App, error) {
	return New(root), nil
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
		timeout = DefaultTimeout
	}

	cltx, cancel := context.WithTimeout(ctx, timeout)
	errCh := make(chan error, 1)

	defer cancel()

	go func() {
		errCh <- f(cltx)
	}()

	select {
	case <-cltx.Done():
		err := cltx.Err()
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("context exceeded set deadline after %v: %w", timeout, err)
		}
		return err
	case err := <-errCh:
		return err
	}
}
