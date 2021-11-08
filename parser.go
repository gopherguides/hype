package hype

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/gopherguides/hype/atomx"
	"github.com/markbates/fsx"
	"golang.org/x/net/html"
)

// Parser will convert HTML documents into, easy to use, nice types.
type Parser struct {
	*fsx.FS
	sync.RWMutex
	customTags   TagMap
	snippetRules map[string]string
	once         sync.Once
}

func (p *Parser) init() {
	p.once.Do(func() {

		if p.snippetRules == nil {
			p.snippetRules = map[string]string{}
		}

		if p.customTags == nil {
			p.customTags = TagMap{}
		}

		p.customTags[atomx.Meta] = func(node *Node) (Tag, error) {
			return NewMeta(node)
		}

		img := func(node *Node) (Tag, error) {
			return NewImage(p.FS, node)
		}

		p.customTags[atomx.Img] = img
		p.customTags[atomx.Image] = img

		p.customTags[atomx.Code] = func(node *Node) (Tag, error) {
			return NewCode(node, p)
		}

		p.customTags[atomx.Body] = func(node *Node) (Tag, error) {
			return NewBody(node)
		}
	})
}

func (p *Parser) SetCustomTag(atom atomx.Atom, fn CustomTagFn) {
	p.init()

	p.Lock()
	p.customTags[atom] = fn
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

	for k, v := range p.customTags {
		p2.customTags[k] = v
	}

	return p2, nil
}

func NewParser(cab fs.FS) (*Parser, error) {
	if cab == nil {
		return nil, fmt.Errorf("cab can not be nil")
	}

	p := &Parser{
		FS:         fsx.NewFS(cab),
		customTags: TagMap{},
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
	p.init()
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return p.NewDocument(node)
}

func (p *Parser) CustomTag(atom atomx.Atom) (CustomTagFn, bool) {
	p.init()

	p.Lock()
	defer p.Unlock()

	fn, ok := p.customTags[atom]
	return fn, ok
}
