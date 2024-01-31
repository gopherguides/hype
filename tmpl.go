package hype

import (
	"bytes"
	"html/template"
	"io"
)

func GoTemplates() PreParseFn {
	fn := func(p *Parser, r io.Reader) (io.Reader, error) {
		if p == nil {
			return nil, ErrIsNil("parser")
		}

		if r == nil {
			return nil, ErrIsNil("reader")
		}

		b, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New("").Parse(string(b))
		if err != nil {
			return nil, err
		}

		bb := &bytes.Buffer{}

		err = tmpl.Execute(bb, p.Vars.Map())
		if err != nil {
			return nil, err
		}

		return bb, nil
	}

	return fn
}
