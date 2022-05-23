package hype

import (
	"io"
)

type PreParser interface {
	PreParse(p *Parser, r io.Reader) (io.Reader, error)
}

type PreParsers []PreParser

func (list PreParsers) PreParse(p *Parser, r io.Reader) (io.Reader, error) {
	var err error

	for _, pp := range list {
		r, err = pp.PreParse(p, r)
		if err != nil {
			return nil, PreParseError{
				Err:       err,
				PreParser: pp,
			}
		}
	}

	return r, nil
}

type PreParseFn func(p *Parser, r io.Reader) (io.Reader, error)

func (fn PreParseFn) PreParse(p *Parser, r io.Reader) (io.Reader, error) {
	return fn(p, r)
}
