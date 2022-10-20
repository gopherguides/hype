package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gopherguides/hype/cmd/hype/cli"
	"github.com/markbates/garlic"
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

	g := &garlic.Garlic{
		Cmd:  app,
		Name: app.Name,
		IO:   app.Stdio(),
	}

	err = g.Main(ctx, pwd, args)
	if err != nil {
		fmt.Fprintln(app.Stderr(), err)
		os.Exit(g.Exit)

		// if errors.Is(err, flag.ErrHelp) || errors.Is(err, cleo.ErrNoCommand) {
		// 	app.Exit(-1, nil)
		// 	return
		// }

		// app.Exit(1, err)
	}
}
