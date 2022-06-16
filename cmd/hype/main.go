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
	"path/filepath"
	"regexp"
	"strconv"

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

	// var fn func([]string) error

	// if len(args) == 0 {
	// 	fn = marked
	// }

	// cmd := args[0]
	// args = args[1:]

	// switch cmd {
	// case "marked", "preview":
	// 	fn = marked
	// 	if len(args) > 0 {
	// 		fn = file
	// 	}
	// // case "vscode":
	// // 	if len(args) == 0 {
	// // 		log.Fatal("missing file")
	// // 	}
	// // 	fn = func() error {
	// // 		return vscode(args[0])
	// // 	}
	// default:
	// 	log.Fatalf("unknown command: %s", cmd)
	// }

	// if err := fn(args); err != nil {
	// 	log.Fatal(err)
	// }

}

func file(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing file")
	}

	name := args[0]

	dir := filepath.Dir(name)
	cab := os.DirFS(dir)
	p := hype.NewParser(cab)

	p.Section = 1

	if sec, err := SectionFromPath(dir); err == nil {
		p.Section = sec
	}

	p.PreParsers = append(p.PreParsers, &Binding{
		Binder: flect.New("book"),
		Ident:  flect.New("chapter"),
	})

	doc, err := p.ParseExecuteFile(context.Background(), filepath.Base(name))
	if err != nil {
		return err
	}

	fmt.Println(doc.String())

	return nil
}

func marked(args []string) error {
	pwd := os.Getenv("MARKED_ORIGIN")

	if len(pwd) > 0 {
		if err := os.Chdir(pwd); err != nil {
			return err
		}
	}

	cab := os.DirFS(".")

	p := hype.NewParser(cab)

	p.Section = 1

	if mp := os.Getenv("MARKED_PATH"); len(mp) > 0 {
		sec, err := SectionFromPath(mp)
		if err != nil {
			return err
		}
		p.Section = sec
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

func SectionFromPath(mp string) (int, error) {
	dir := filepath.Dir(mp)
	base := filepath.Base(dir)
	rx, err := regexp.Compile(`^(\d+)-.+`)
	if err != nil {
		return 0, err
	}

	match := rx.FindStringSubmatch(base)
	if len(match) < 2 {
		return 0, fmt.Errorf("could not find section: %q", mp)
	}

	sec, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}
	return sec, nil
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

// func vscode(name string) error {

// 	dir := filepath.Dir(name)
// 	if len(dir) > 0 {
// 		if err := os.Chdir(dir); err != nil {
// 			return err
// 		}
// 	}
// 	cab := os.DirFS(".")

// 	p := hype.NewParser(cab)

// 	p.Section = 1
// 	p.PreParsers = append(p.PreParsers, &Binding{
// 		Binder: flect.New("book"),
// 		Ident:  flect.New("chapter"),
// 	})

// 	doc, err := p.ParseExecuteFile(context.Background(), filepath.Base(name))
// 	if err != nil {
// 		return err
// 	}

// 	body, err := doc.Body()
// 	if err != nil {
// 		return err
// 	}

// 	return json.NewEncoder(os.Stdout).Encode(map[string]any{
// 		"body":  body.Nodes.String(),
// 		"title": doc.Title,
// 	})
// }
