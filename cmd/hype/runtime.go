package main

import (
	"flag"
	"io"
	"io/fs"
	"os"

	"github.com/markbates/fsx"
)

type Runtime struct {
	*flag.FlagSet
	Args       []string
	Cab        fs.FS
	Stderr     io.Writer
	Stdin      io.Reader
	Stdout     io.Writer
	IndentJSON bool
}

func (rt *Runtime) Parse(args []string) error {
	if err := rt.FlagSet.Parse(args); err != nil {
		return err
	}

	rt.Args = rt.FlagSet.Args()
	return nil
}

func runtime(args []string) (*Runtime, error) {
	rt := &Runtime{
		Args:    args,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		FlagSet: flag.NewFlagSet("", flag.ExitOnError),
	}
	rt.SetOutput(rt.Stdout)
	rt.BoolVar(&rt.IndentJSON, "i", false, "indent json")

	if err := rt.Parse(args); err != nil {
		return nil, err
	}

	pwd, _ := os.Getwd()
	cab, err := fsx.DirFS(pwd)
	if err != nil {
		return nil, err
	}

	rt.Cab = cab

	return rt, nil
}
