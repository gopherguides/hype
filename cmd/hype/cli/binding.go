package cli

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
)

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
