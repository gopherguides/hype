package hype

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gopherguides/hype/mdx"
)

func Markdown() PreParseFn {
	fn := func(p *Parser, r io.Reader) (io.Reader, error) {
		if p == nil {
			return nil, ErrIsNil("parser")
		}

		if r == nil {
			return nil, ErrIsNil("reader")
		}

		md := mdx.New()
		md.DisablePages = p.DisablePages

		b, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("markdown: %w", err)
		}
		b = bytes.ReplaceAll(b, []byte("\\n"), []byte("  \n"))

		b, err = md.Parse(b)
		if err != nil {
			return nil, fmt.Errorf("markdown: %w", err)
		}

		b = bytes.ReplaceAll(b, []byte("&rsquo;"), []byte("'"))
		b = bytes.ReplaceAll(b, []byte("&ldquo;"), []byte("\""))
		b = bytes.ReplaceAll(b, []byte("&rdquo;"), []byte("\""))

		return bytes.NewReader(b), nil
	}

	return fn
}
