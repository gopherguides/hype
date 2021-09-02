package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	case "toc":
		rt.Args = args[1:]
		return toc(rt)
	}

	doc, err := parseFile(rt)
	if err != nil {
		rt.Usage()
		return err
	}

	return print(rt.Stdout, doc)
}

func print(w io.Writer, doc *hype.Document) error {
	p := hype.NewPrinter(w)
	p.SetTransformer(func(tag hype.Tag) (hype.Tag, error) {

		inc, ok := tag.(*hype.Include)
		if !ok {
			return tag, nil
		}

		src := inc.Src()
		// base := filepath.Base(src)
		dir := filepath.Dir(src)

		for _, code := range inc.Children.AllType(&hype.Image{}) {
			sc, ok := code.(*hype.Image)
			if !ok {
				continue
			}
			x := sc.Src()
			if strings.HasPrefix(x, "http") {
				continue
			}
			x = filepath.Join(dir, x)
			sc.Set("src", x)
		}

		return inc, nil

	})

	return p.Print(doc.Children...)
}

func toc(rt *Runtime) error {
	// doc, err := parseFile(rt)
	// if err != nil {
	// 	rt.Usage()
	// 	return err
	// }

	// err = hype.Print(doc.Children, rt.Stdout, func(w io.Writer, t hype.Tag) error {
	// 	text := t.GetChildren().String()
	// 	switch t.Atom() {
	// 	case atom.H1:
	// 		fmt.Fprintln(w, text)
	// 		return nil
	// 	case atom.H2:
	// 		fmt.Fprintf(w, "  %s\n", text)
	// 	case atom.H3:
	// 		fmt.Fprintf(w, "    %s\n", text)
	// 	case atom.H4:
	// 		fmt.Fprintf(w, "      %s\n", text)
	// 	case atom.H5:
	// 		fmt.Fprintf(w, "        %s\n", text)
	// 	case atom.H6:
	// 		fmt.Fprintf(w, "          %s\n", text)
	// 	}
	// 	return nil
	// })

	// return err

	return nil
}

func jsonCmd(rt *Runtime) error {
	doc, err := parseFile(rt)
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

	return print(rt.Stdout, doc)
}

func parseFile(rt *Runtime) (*hype.Document, error) {
	args := rt.Args

	if len(args) == 0 {
		return nil, fmt.Errorf("missing file name")
	}

	fn := args[0]
	base := filepath.Base(fn)
	dir := filepath.Dir(fn)
	rt.Cab = os.DirFS(dir)

	p, err := hype.NewParser(rt.Cab)
	if err != nil {
		return nil, err
	}

	return p.ParseFile(base)
}
