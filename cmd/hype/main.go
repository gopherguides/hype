package main

import (
	"context"
	_ "embed"
	"errors"
	"log"
	"os"
	"os/exec"
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

	garlic := &garlic.Garlic{
		Cmd:  app,
		Name: app.Name,
		IO:   app.Stdio(),
	}

	err = garlic.Main(ctx, pwd, args)
	if err != nil {
		var ex *exec.ExitError
		if ok := errors.As(err, &ex); ok {
			os.Exit(ex.ExitCode())
		}

		os.Exit(1)
		// if errors.Is(err, flag.ErrHelp) || errors.Is(err, cleo.ErrNoCommand) {
		// 	app.Exit(-1, nil)
		// 	return
		// }

		// app.Exit(1, err)
	}
}
