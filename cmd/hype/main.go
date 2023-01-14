package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/gopherguides/hype/cmd/hype/cli"
	"github.com/markbates/cleo"
	"github.com/markbates/garlic"
)

func main() {
	args := os.Args[1:]

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	xp := os.Getenv("PATH")
	paths := []string{
		xp,
		"/opt/homebrew/bin",
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
		"/usr/sbin",
		"/sbin",
	}

	xp = strings.Join(paths, ":")
	os.Setenv("PATH", xp)

	if mp := os.Getenv("MARKED_PATH"); len(mp) > 0 {
		pwd = filepath.Dir(mp)
	}

	app := cli.New(pwd)
	app.Name = "hype"

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	clove := &garlic.Garlic{
		Name: "hype",
		Cmd:  app,
		FS:   os.DirFS(pwd),
	}

	err = clove.Main(ctx, pwd, args)
	if err != nil {
		err = fmt.Errorf("%w: PATH: %q", err, xp)
		cleo.Exit(app, 1, err)
		os.Exit(1)
	}
}
