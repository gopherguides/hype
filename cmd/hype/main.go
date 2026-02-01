package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/signal"

	"github.com/gopherguides/hype/cmd/hype/cli"
	"github.com/markbates/cleo"
	"github.com/markbates/garlic"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	args := os.Args[1:]

	pwd, err := cli.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	app := cli.New(pwd, cli.VersionInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
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
		cleo.Exit(app, 1, err)
		os.Exit(1)
	}
}
