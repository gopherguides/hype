package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gopherguides/hype/cmd/hype/cli"
	"github.com/markbates/cleo"
)

func main() {
	args := os.Args[1:]

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if mp := os.Getenv("MARKED_PATH"); len(mp) > 0 {
		pwd = filepath.Dir(mp)
	}

	app := cli.New(pwd)
	app.Name = "hype"

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err = app.Main(ctx, pwd, args)
	if err != nil {
		cleo.Exit(app, 1, err)
		os.Exit(1)
	}
}
