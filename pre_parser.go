package hype

import (
	"bytes"
	"io"
)

type PreParser interface {
	PreParse(p *Parser, r io.Reader) (io.Reader, error)
}

type PreParsers []PreParser

func (list PreParsers) PreParse(p *Parser, r io.Reader) (io.Reader, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(contents)

	for _, pp := range list {
		r, err = pp.PreParse(p, r)
		if err != nil {
			return nil, PreParseError{
				Contents:  contents,
				Err:       err,
				Filename:  p.Filename,
				PreParser: pp,
				Root:      p.Root,
			}
		}
	}

	return r, nil
}

type PreParseFn func(p *Parser, r io.Reader) (io.Reader, error)

func (fn PreParseFn) PreParse(p *Parser, r io.Reader) (io.Reader, error) {
	return fn(p, r)
}
