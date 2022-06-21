package hype

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sync"

	"golang.org/x/net/html"
)

type ParseElementFn func(p *Parser, el *Element) (Nodes, error)

type Parser struct {
	fs.FS

	Root         string
	DisablePages bool
	NodeParsers  map[Atom]ParseElementFn
	PreParsers   PreParsers
	Snippets     *Snippets
	Section      int

	fileName string
	mu       sync.RWMutex
}

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

func (p *Parser) Sub(dir string) (*Parser, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	cab, err := fs.Sub(p.FS, dir)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	p2 := &Parser{
		FS:          cab,
		Root:        filepath.Join(p.Root, dir),
		PreParsers:  p.PreParsers,
		NodeParsers: p.NodeParsers,
	}

	return p2, nil
}

func NewParser(cab fs.FS) *Parser {
	return &Parser{
		FS:          cab,
		NodeParsers: DefaultElements(),
		PreParsers:  PreParsers{Markdown()},
		Section:     1,
		Snippets:    &Snippets{},
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

	return parserErr{
		p:   p,
		err: err,
	}
}

type parserErr struct {
	p   *Parser
	err error
}

func (pe parserErr) Error() string {
	if pe.err == nil {
		return ""
	}

	err := pe.err

	p := pe.p

	if p == nil {
		return err.Error()
	}

	if len(p.fileName) > 0 {
		err = fmt.Errorf("file: %q: %s", p.fileName, err)
	}

	if len(p.Root) > 0 {
		err = fmt.Errorf("root: %q: %s", p.Root, err)
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
