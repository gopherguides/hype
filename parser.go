package hype

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/markbates/fsx"
	"golang.org/x/net/html"
)

type CustomTagFn func(node *Node) (Tag, error)

// Parser will convert HTML documents into, easy to use, nice types.
type Parser struct {
	*fsx.FS
	*sync.RWMutex
	customTags   map[string]CustomTagFn
	snippetRules map[string]string
}

func (p *Parser) SetCustomTag(name string, fn CustomTagFn) {
	p.Lock()
	if p.customTags == nil {
		p.customTags = map[string]CustomTagFn{}
	}
	p.customTags[name] = fn
	p.Unlock()
}

func (p *Parser) SubParser(path string) (*Parser, error) {
	p.Lock()
	defer p.Unlock()

	cab, err := p.Sub(path)
	if err != nil {
		return nil, err
	}

	p2, err := NewParser(cab)
	if err != nil {
		return nil, err
	}

	for k, v := range p.snippetRules {
		p2.snippetRules[k] = v
	}

	return p2, nil
}

func NewParser(cab fs.FS) (*Parser, error) {
	if cab == nil {
		return nil, fmt.Errorf("cab can not be nil")
	}

	p := &Parser{
		FS:         fsx.NewFS(cab),
		RWMutex:    &sync.RWMutex{},
		customTags: map[string]CustomTagFn{},
		snippetRules: map[string]string{
			".html": "<!-- %s -->",
			".go":   "// %s",
		},
	}

	return p, nil
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
