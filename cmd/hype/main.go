package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
	"github.com/markbates/fsx"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	sw := cleo.NewRouter()
	sw.SetDefault(defaultSW(sw))

	rt, err := cleo.NewRuntime("hype", os.Args[1:])
	if err != nil {
		return err
	}

	return sw.Switch(rt)
}

func defaultSW(sw *cleo.Router) cleo.HandlerFn {
	return func(rt *cleo.Runtime) error {
		var err error
		pwd := os.Getenv("MARKED_ORIGIN")
		if len(pwd) == 0 {
			pwd, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		cab, err := fsx.DirFS(pwd)
		if err != nil {
			return err
		}
		rt.Cab = cab

		p, err := hype.NewParser(rt.Cab)
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(rt.Stdin)
		if err != nil {
			return err
		}

		doc, err := p.ParseMD(b)
		if err != nil {
			return err
		}

		fmt.Fprintln(rt.Stdout, doc.String())
		return nil
	}
}
