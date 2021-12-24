package hype

import (
	"bytes"
	"io"

	"github.com/gopherguides/hype/mdx"
)

// ParseMD parses the given markdown into a Document.
func (p *Parser) ParseMD(src []byte) (*Document, error) {
	mdp := mdx.New()

	src, err := mdp.Parse(src)
	if err != nil {
		return nil, err
	}

	r := io.NopCloser(bytes.NewReader(src))

	defer r.Close()

	doc, err := p.ParseReader(r)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
