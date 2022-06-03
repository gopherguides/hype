package hype

import (
	"fmt"
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

type PreParseError struct {
	Err       error
	PreParser PreParser
}

func (e PreParseError) Error() string {
	return fmt.Sprintf("pre parse error: [%T]: %v", e.PreParser, e.Err)
}
