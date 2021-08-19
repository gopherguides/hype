package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gopherguides/hype"
	"github.com/markbates/fsx"
)

func main() {
	rt, err := runtime(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if err := run(rt); err != nil {
		log.Fatal(err)
	}
}

func run(rt *Runtime) error {
	args := rt.Args

	if len(args) == 0 {
		return stream(rt)
	}

	switch args[0] {
	case "json":
		rt.Args = args[1:]
		return jsonCmd(rt)
	}

	return fmt.Errorf("unknown arguments %q", args)
}

func jsonCmd(rt *Runtime) error {
	args := rt.Args

	if len(args) == 0 {
		return fmt.Errorf("missing file name")
	}

	fn := args[0]
	base := filepath.Base(fn)
	dir := filepath.Dir(fn)
	rt.Cab = os.DirFS(dir)

	p, err := hype.NewParser(rt.Cab)
	if err != nil {
		return err
	}

	doc, err := p.ParseFile(base)
	if err != nil {
		rt.Usage()
		return err
	}

	enc := json.NewEncoder(rt.Stdout)

	if rt.IndentJSON {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(doc); err != nil {
		return err
	}

	return nil
}

func stream(rt *Runtime) error {
	pwd := os.Getenv("MARKED_ORIGIN")
	if len(pwd) == 0 {
		pwd, _ = os.Getwd()
	}

	cab, err := fsx.DirFS(pwd)
	if err != nil {
		return err
	}

	p, err := hype.NewParser(cab)
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

	fmt.Fprint(rt.Stdout, doc)
	return nil
}
