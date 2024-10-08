package hype

import (
	"bytes"
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
	fs.FS

	DisablePages bool
	DocIDGen     func() (string, error) // default: uuid.NewV4().String()
	Filename     string                 // only set when Parser.ParseFile() is used
	NodeParsers  map[Atom]ParseElementFn
	NowFn        func() time.Time // default: time.Now()
	PreParsers   PreParsers
	Root         string
	Section      int
	Vars         syncx.Map[string, any]
	Contents     []byte // a copy of the contents being parsed - set just before parsing

	mu sync.RWMutex
}

func (p *Parser) MarshalJSON() ([]byte, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	x := struct {
		Type         string         `json:"type,omitempty"`
		Root         string         `json:"root,omitempty"`
		DisablePages bool           `json:"disable_pages,omitempty"`
		Section      int            `json:"section,omitempty"`
		Vars         map[string]any `json:"vars,omitempty"`
		Contents     string         `json:"contents,omitempty"`
	}{
		Type:         toType(p),
		Root:         p.Root,
		DisablePages: p.DisablePages,
		Section:      p.Section,
		Vars:         p.Vars.Map(),
		Contents:     string(p.Contents),
	}

	return json.MarshalIndent(x, "", "  ")

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
func (p *Parser) ParseFile(name string) (doc *Document, err error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	defer func() {
		err = p.ensureParseError(err)
	}()

	p.mu.Lock()
	p.Filename = name
	p.mu.Unlock()

	f, err := p.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	doc, err = p.Parse(f)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (p *Parser) Parse(r io.Reader) (doc *Document, err error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	defer func() {
		err = p.ensureParseError(err)
	}()

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	p.Contents = b

	r = bytes.NewReader(b)

	// pre parse
	r, err = p.PreParsers.PreParse(p, r)
	if err != nil {
		return nil, err
	}

	// parse
	hdoc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	node, err := p.ParseHTMLNode(hdoc, nil)
	if err != nil {
		return nil, err
	}

	doc, err = p.newDoc()
	if err != nil {
		return nil, err
	}

	doc.Nodes = Nodes{node}
	if len(doc.Title) == 0 {
		doc.Title = FindTitle(doc.Nodes)
	}

	// post parse
	err = doc.Nodes.PostParse(p, doc, err)
	if err != nil {
		return nil, err
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

func (p *Parser) ParseFolder(name string) (doc Documents, err error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	defer func() {
		if err == nil {
			return
		}

		if _, ok := err.(ParseError); ok {
			return
		}

		err = ParseError{
			Err:      err,
			Filename: name,
			Root:     p.Root,
		}
	}()

	var docs Documents
	var wg errgroup.Group
	var mu sync.Mutex

	whole, err := binding.WholeFromPath(p.FS, name, "book", "chapter")
	if err != nil && !errors.Is(err, binding.ErrPath("")) {
		return nil, err
	}

	err = fs.WalkDir(p.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)

		if base != "hype.md" {
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
	case html.DocumentNode, html.ElementNode:
		n, err := p.element(node, parent)
		if err != nil {
			return nil, err
		}
		return n, nil
	case html.TextNode:
		return Text(node.Data), nil
	case html.ErrorNode:
		return nil, ParseError{
			Err:      fmt.Errorf("error node: %+v", node),
			Filename: p.Filename,
			Root:     p.Root,
		}
	}

	return nil, ParseError{
		Err:      fmt.Errorf("unknown node: %+v", node),
		Filename: p.Filename,
		Root:     p.Root,
	}
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
		return nil, err
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
		Filename:   p.Filename,
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
		FS:        p.FS,
		Filename:  p.Filename,
		ID:        id,
		Parser:    p,
		Root:      p.Root,
		SectionID: p.Section,
		Snippets:  Snippets{},
	}

	return doc, nil
}

func (p *Parser) ensureParseError(err error) error {
	if p == nil {
		return ErrIsNil("parser")
	}

	if err == nil {
		return nil
	}

	if _, ok := err.(ParseError); ok {
		return err
	}

	return ParseError{
		Err:      err,
		Filename: p.Filename,
		Root:     p.Root,
		Contents: p.Contents,
	}
}
