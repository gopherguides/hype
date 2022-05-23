package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
)

func main() {
	args := os.Args[1:]

	var fn func() error
	if len(args) == 0 {
		fn = marked
	}

	cmd := args[0]
	args = args[1:]

	switch cmd {
	case "marked", "preview":
		fn = marked
	case "vscode":
		if len(args) == 0 {
			log.Fatal("missing file")
		}
		fn = func() error {
			return vscode(args[0])
		}
	default:
		log.Fatalf("unknown command: %s", cmd)
	}

	if err := fn(); err != nil {
		log.Fatal(err)
	}

}

func marked() error {
	pwd := os.Getenv("MARKED_ORIGIN")

	if len(pwd) > 0 {
		os.Chdir(pwd)
	}

	cab := os.DirFS(".")

	p := hype.NewParser(cab)

	p.Section = 1

	if mp := os.Getenv("MARKED_PATH"); len(mp) > 0 {
		dir := filepath.Dir(mp)
		base := filepath.Base(dir)
		rx, err := regexp.Compile(`^(\d+)-.+`)
		if err != nil {
			return err
		}

		match := rx.FindStringSubmatch(base)
		if len(match) >= 2 {
			sec, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}
			p.Section = sec
		}
	}
	p.PreParsers = append(p.PreParsers, &Binding{
		Binder: flect.New("book"),
		Ident:  flect.New("chapter"),
	})

	doc, err := p.ParseExecute(context.Background(), os.Stdin)
	if err != nil {
		return err
	}

	fmt.Println(doc.String())

	return nil
}

type Binding struct {
	Binder flect.Ident // book
	Ident  flect.Ident // chapter
}

func (bind *Binding) String() string {
	if bind == nil {
		return ""
	}

	return bind.Ident.String()
}

func (bind *Binding) PreParse(p *hype.Parser, r io.Reader) (io.Reader, error) {
	if r == nil {
		return r, nil
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	in := string(b)

	tmpl, err := template.New("").Parse(in)
	if err != nil {
		return nil, fmt.Errorf("parse: %w: %s", err, in)
	}

	bb := &bytes.Buffer{}

	err = tmpl.Execute(bb, bind)
	if err != nil {
		return nil, fmt.Errorf("execute: %w: %s", err, in)
	}

	return bb, nil
}

func vscode(name string) error {

	dir := filepath.Dir(name)
	if len(dir) > 0 {
		os.Chdir(dir)
	}
	cab := os.DirFS(".")

	p := hype.NewParser(cab)

	p.Section = 1
	p.PreParsers = append(p.PreParsers, &Binding{
		Binder: flect.New("book"),
		Ident:  flect.New("chapter"),
	})

	doc, err := p.ParseExecuteFile(context.Background(), filepath.Base(name))
	if err != nil {
		return err
	}

	body, err := doc.Body()
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(map[string]any{
		"body":  body.Nodes.String(),
		"title": doc.Title,
	})
}
