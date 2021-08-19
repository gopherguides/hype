package hype

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/markbates/fsx"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
)

// Parser will convert HTML documents into, easy to use, nice types.
type Parser struct {
	*fsx.FS
	*sync.RWMutex
	snippetRules map[string]string
}

func NewParser(cab fs.FS) (*Parser, error) {
	if cab == nil {
		return nil, fmt.Errorf("cab can not be nil")
	}

	p := &Parser{
		FS:      fsx.NewFS(cab),
		RWMutex: &sync.RWMutex{},
		snippetRules: map[string]string{
			".html": "<!-- %s -->",
			".go":   "// %s",
		},
	}

	return p, nil
}

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
	r := p.markdown(src)
	defer r.Close()

	doc, err := p.ParseReader(r)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// ParseFile will parse the requested file and return a Document.
// The file MUST be present in the parser's FS.
func (p *Parser) ParseFile(name string) (*Document, error) {
	var r io.ReadCloser
	r, err := p.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if filepath.Ext(name) == ".md" {
		src, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return p.ParseMD(src)
	}

	doc, err := p.ParseReader(r)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (p *Parser) ParseReader(r io.ReadCloser) (*Document, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return p.NewDocument(node)
}
