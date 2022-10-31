package hype

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/gopherguides/hype/binding"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

type ParseElementFn func(p *Parser, el *Element) (Nodes, error)

type Parser struct {
	fs.FS

	Root         string
	DisablePages bool
	NodeParsers  map[Atom]ParseElementFn
	PreParsers   PreParsers
	Snippets     Snippets
	Section      int
	NowFn        func() time.Time // default: time.Now()

	fileName string
	mu       sync.RWMutex
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
		return nil, p.wrapErr(err)
	}
	defer f.Close()

	doc, err := p.Parse(f)
	if err != nil {
		return nil, p.wrapErr(err)
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
		return nil, p.wrapErr(err)
	}

	// parse
	hdoc, err := html.Parse(r)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	node, err := p.ParseHTMLNode(hdoc, nil)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	doc := p.newDoc()
	doc.Nodes = Nodes{node}
	if len(doc.Title) == 0 {
		doc.Title = FindTitle(doc.Nodes)
	}

	// post parse
	err = doc.Nodes.PostParse(p, doc, err)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	return doc, nil
}

func (p *Parser) ParseExecuteFile(ctx context.Context, name string) (*Document, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	doc, err := p.ParseFile(name)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	err = doc.Execute(ctx)
	if err != nil {
		return nil, p.wrapErr(err)
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
		return nil, p.wrapErr(err)
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
				return fmt.Errorf("error parsing: %q: %w", path, err)
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
		return nil, p.wrapErr(err)
	}

	if err := wg.Wait(); err != nil {
		return nil, p.wrapErr(err)
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
		return nil, p.wrapErr(err)
	}

	err = docs.Execute(ctx)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	return docs, nil
}

func (p *Parser) ParseExecuteFragment(ctx context.Context, r io.Reader) (Nodes, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	nodes, err := p.ParseFragment(r)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	doc := p.newDoc()
	doc.Nodes = nodes

	err = doc.Execute(ctx)

	if err != nil {
		return nil, p.wrapErr(err)
	}

	return nodes, nil
}

func (p *Parser) ParseExecute(ctx context.Context, r io.Reader) (*Document, error) {
	doc, err := p.Parse(r)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	err = doc.Execute(ctx)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	return doc, nil
}

func (p *Parser) ParseFragment(r io.Reader) (Nodes, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	doc, err := p.Parse(r)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	pages := ByType[*Page](doc.Nodes)

	if len(pages) > 0 {
		return pages[0].Nodes, nil
	}

	body, err := doc.Body()
	if err != nil {
		return nil, p.wrapErr(err)
	}

	return body.Nodes, nil
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
			return nil, p.wrapErr(err)
		}
		return n, nil
	case html.ErrorNode:
		return nil, p.wrapErr(fmt.Errorf(node.Data))
	case html.TextNode:
		return Text(node.Data), nil
	}

	return nil, p.wrapErr(fmt.Errorf("unknown node type %v", node.Data))
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
		return nil, p.wrapErr(err)
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
		FileName:   p.fileName,
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
		PreParsers:  PreParsers{Markdown()},
		Section:     1,
	}
}

func (p *Parser) newDoc() *Document {
	doc := &Document{
		FS:        p.FS,
		Parser:    p,
		Root:      p.Root,
		SectionID: p.Section,
		Snippets:  p.Snippets,
	}

	return doc
}

func (p *Parser) wrapErr(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(parserErr); ok {
		return err
	}

	if p == nil {
		return err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	return parserErr{
		fileName: p.fileName,
		root:     p.Root,
		err:      err,
	}
}

type parserErr struct {
	err      error
	fileName string
	root     string
}

func (pe parserErr) Error() string {
	if pe.err == nil {
		return ""
	}

	err := pe.err

	if len(pe.fileName) > 0 {
		err = fmt.Errorf("file: %q: %s", pe.fileName, err)
	}

	if len(pe.root) > 0 {
		err = fmt.Errorf("root: %q: %s", pe.root, err)
	}

	return err.Error()
}

func (pe parserErr) Unwrap() error {
	type Unwrapper interface {
		Unwrap() error
	}

	if _, ok := pe.err.(Unwrapper); ok {
		return errors.Unwrap(pe.err)
	}

	return pe.err
}

func (pe parserErr) As(target any) bool {
	ex, ok := target.(*parserErr)
	if !ok {
		return errors.As(pe.err, target)
	}

	(*ex) = pe
	return true
}

func (pe parserErr) Is(target error) bool {
	if _, ok := target.(parserErr); ok {
		return true
	}

	return errors.Is(pe.err, target)
}
