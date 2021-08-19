package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gopherguides/hype"

	"github.com/markbates/fsx"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	pwd := os.Getenv("MARKED_ORIGIN")
	// fmt.Println("TODO >> main.go:27 pwd", pwd)
	// panic(pwd)
	// if len(pwd) == 0 {
	// 	pwd, _ = os.Getwd()
	// }

	cab, err := fsx.DirFS(pwd)
	if err != nil {
		return err
	}

	p, err := hype.NewParser(cab)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	doc, err := p.ParseMD(b)
	if err != nil {
		return err
	}

	fmt.Println(doc.String())
	return nil
}
