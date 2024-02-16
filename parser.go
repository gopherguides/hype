package hype

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gopherguides/hype/binding"
	"github.com/markbates/syncx"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

type ParseElementFn func(p *Parser, el *Element) (Nodes, error)

type Parser struct {
	fs.FS `json:"-"`

	Root         string                  `json:"root,omitempty"`
	DisablePages bool                    `json:"disable_pages,omitempty"`
	NodeParsers  map[Atom]ParseElementFn `json:"-"`
	PreParsers   PreParsers              `json:"-"`
	Snippets     Snippets                `json:"snippets,omitempty"`
	Section      int                     `json:"section,omitempty"`
	NowFn        func() time.Time        `json:"-"` // default: time.Now()
	DocIDGen     func() (string, error)  `json:"-"` // default: uuid.NewV4().String()
	Vars         syncx.Map[string, any]  `json:"vars,omitempty"`

	fileName string
	mu       sync.RWMutex
}

func (p *Parser) MarshalJSON() ([]byte, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	// p.Vars.Map()

	x := struct {
		Type         string         `json:"type,omitempty"`
		Root         string         `json:"root,omitempty"`
		DisablePages bool           `json:"disable_pages,omitempty"`
		Section      int            `json:"section,omitempty"`
		Snippets     Snippets       `json:"snippets,omitempty"`
		Vars         map[string]any `json:"vars,omitempty"`
	}{
		Type:         fmt.Sprintf("%T", p),
		Root:         p.Root,
		DisablePages: p.DisablePages,
		Section:      p.Section,
		Snippets:     p.Snippets,
		Vars:         p.Vars.Map(),
	}

	return json.Marshal(x)

}

func (p *Parser) Now() time.Time {
	if p == nil || p.NowFn == nil {
		return time.Now()
	}

	return p.NowFn()
}

// ParseFile parses the given file from Parser.FS.
// If successful a *Document is returned. The returned
// *Document is NOT yet executed.
func (p *Parser) ParseFile(name string) (*Document, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	p.mu.Lock()
	p.fileName = name
	p.mu.Unlock()

	f, err := p.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	doc, err := p.Parse(f)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (p *Parser) Parse(r io.Reader) (*Document, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	// pre parse
	r, err := p.PreParsers.PreParse(p, r)
	if err != nil {
		return nil, p.newError(err)
	}

	// parse
	hdoc, err := html.Parse(r)
	if err != nil {
		return nil, p.newError(err)
	}

	node, err := p.ParseHTMLNode(hdoc, nil)
	if err != nil {
		return nil, p.newError(err)
	}

	doc, err := p.newDoc()
	if err != nil {
		return nil, p.newError(err)
	}

	doc.Nodes = Nodes{node}
	if len(doc.Title) == 0 {
		doc.Title = FindTitle(doc.Nodes)
	}

	// post parse
	err = doc.Nodes.PostParse(p, doc, err)
	if err != nil {
		return nil, p.newError(err)
	}

	return doc, nil
}

func (p *Parser) ParseExecuteFile(ctx context.Context, name string) (*Document, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	doc, err := p.ParseFile(name)
	if err != nil {
		return nil, err
	}

	err = doc.Execute(ctx)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (p *Parser) ParseFolder(name string) (Documents, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	var docs Documents
	var wg errgroup.Group
	var mu sync.Mutex

	whole, err := binding.WholeFromPath(p.FS, name, "book", "chapter")
	if err != nil && !errors.Is(err, binding.ErrPath("")) {
		pe := p.newError(err)
		pe.Filename = name
		return nil, pe
	}

	err = fs.WalkDir(p.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)

		if base != "module.md" {
			return nil
		}

		wg.Go(func() error {
			dir := filepath.Dir(path)

			var part binding.Part
			for _, v := range whole.Parts {
				db := filepath.Base(dir)
				pb := filepath.Base(v.Path)
				if db == pb {
					part = v
					break
				}
			}

			p, err := p.Sub(dir)
			if err != nil {
				return fmt.Errorf("error getting sub fs: %q: %w", dir, err)
			}

			p.Section = part.Number

			doc, err := p.ParseFile(base)
			if err != nil {
				return err
			}

			mu.Lock()
			docs = append(docs, doc)
			mu.Unlock()

			return nil
		})

		return filepath.SkipDir
	})

	if err != nil {
		err = fmt.Errorf("error walking: %q: %w", name, err)
		if err != nil && !errors.Is(err, binding.ErrPath("")) {
			return nil, err
		}
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].SectionID < docs[j].SectionID
	})

	return docs, nil
}

func (p *Parser) ParseExecuteFolder(ctx context.Context, name string) (Documents, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	docs, err := p.ParseFolder(name)
	if err != nil {
		return nil, err
	}

	err = docs.Execute(ctx)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (p *Parser) ParseExecuteFragment(ctx context.Context, r io.Reader) (Nodes, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	nodes, err := p.ParseFragment(r)
	if err != nil {
		return nil, err
	}

	doc, err := p.newDoc()
	if err != nil {
		return nil, err
	}
	doc.Nodes = nodes

	err = doc.Execute(ctx)

	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (p *Parser) ParseExecute(ctx context.Context, r io.Reader) (*Document, error) {
	doc, err := p.Parse(r)
	if err != nil {
		return nil, err
	}

	err = doc.Execute(ctx)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (p *Parser) ParseFragment(r io.Reader) (Nodes, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	doc, err := p.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := doc.Nodes

	pages := ByType[*Page](nodes)

	if len(pages) > 0 {
		nodes = pages[0].Nodes
	}

	pg, ok := FirstByType[*Paragraph](nodes)
	if ok {
		return pg.Nodes, nil
	}

	return nodes, nil
}

func (p *Parser) ParseHTMLNode(node *html.Node, parent Node) (Node, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	if node == nil {
		return nil, ErrIsNil("node")
	}

	switch node.Type {
	case html.CommentNode:
		return Comment(node.Data), nil
	case html.DoctypeNode:
		// return p.NewDocType(node)
	case html.DocumentNode, html.ElementNode:
		n, err := p.element(node, parent)
		if err != nil {
			return nil, err
		}
		return n, nil
	case html.ErrorNode:
		return nil, fmt.Errorf("error node: %v", node.Data)
	case html.TextNode:
		return Text(node.Data), nil
	}

	return nil, fmt.Errorf("unknown node type %v", node.Data)
}

func (p *Parser) Sub(dir string) (*Parser, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	p2 := &Parser{
		FS:          p.FS,
		Root:        filepath.Join(p.Root, dir),
		PreParsers:  p.PreParsers,
		NodeParsers: p.NodeParsers,
	}

	if len(dir) == 0 || dir == "." {
		return p2, nil
	}

	cab, err := fs.Sub(p.FS, dir)
	if err != nil {
		return nil, p.newError(err)
	}

	p2.FS = cab

	return p2, nil
}

func (p *Parser) element(node *html.Node, parent Node) (Node, error) {
	ats, err := ConvertHTMLAttrs(node.Attr)
	if err != nil {
		return nil, err
	}

	el := &Element{
		Attributes: ats,
		HTMLNode:   node,
		Parent:     parent,
		Filename:   p.fileName,
	}

	var nodes Nodes
	c := node.FirstChild

	for c != nil {
		kid, err := p.ParseHTMLNode(c, el)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, kid)
		c = c.NextSibling
	}

	el.Nodes = nodes

	fn, ok := p.NodeParsers[el.Atom()]
	if ok {
		return fn(p, el)
	}

	return el, nil
}

// NewParser returns a fully initialized Parser.
// This includes the Markdown pre-parser and the
// default node parsers.
func NewParser(cab fs.FS) *Parser {
	return &Parser{
		FS:          cab,
		NodeParsers: DefaultElements(),
		PreParsers:  PreParsers{VarProcessor(), GoTemplates(), Markdown()},
		Section:     1,
		Vars:        syncx.Map[string, any]{},
	}
}

func (p *Parser) newDoc() (*Document, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.DocIDGen == nil {
		p.DocIDGen = func() (string, error) {
			id, err := uuid.NewV4()
			if err != nil {
				return "", err
			}
			return id.String(), nil
		}
	}

	id, err := p.DocIDGen()
	if err != nil {
		return nil, err
	}
	doc := &Document{
		ID:        id,
		FS:        p.FS,
		Parser:    p,
		Root:      p.Root,
		SectionID: p.Section,
		Snippets:  p.Snippets,
		Filename:  p.fileName,
	}

	return doc, nil
}

func (p *Parser) newError(err error) ParseError {
	return ParseError{
		HypeError: HypeError{
			Err:      err,
			Root:     p.Root,
			Filename: p.fileName,
		},
	}
}
