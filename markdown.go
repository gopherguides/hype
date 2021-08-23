package hype

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/russross/blackfriday"
)

func (p *Parser) markdown(src []byte) io.ReadCloser {
	const extensions = blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_TABLES

	r := blackfriday.HtmlRenderer(0, "", "")
	src = blackfriday.Markdown(src, r, extensions)
	return io.NopCloser(bytes.NewReader(src))
}

func (p *Parser) ParseMD(src []byte) (*Document, error) {
	var err error

	r := p.markdown(src)
	if !p.IgnoreMDPages {
		r, err = p.mdPages(r)
		if err != nil {
			return nil, err
		}
	}

	defer r.Close()

	doc, err := p.ParseReader(r)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (p *Parser) mdPages(r io.ReadCloser) (io.ReadCloser, error) {
	defer r.Close()
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	bb := &bytes.Buffer{}

	fmt.Fprintln(bb, `<page number="1">`)

	ind := 1

	for _, line := range bytes.Split(src, []byte("\n")) {
		if !bytes.HasPrefix(line, []byte("<hr")) {
			fmt.Fprintln(bb, string(line))
			continue
		}
		ind++
		fmt.Fprintln(bb, `</page>`)
		fmt.Fprintf(bb, "<page number=\"%d\">", ind)

	}

	fmt.Fprintln(bb, `</page>`)

	return io.NopCloser(bb), nil
}
