package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/cmd/hype/cli"
)

func main() {
	args := os.Args[1:]

	app := cli.New()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err = app.Main(ctx, pwd, args)
	if err != nil {
		log.Fatal(err)
	}
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
