package hype

import (
	"context"
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

	mu sync.RWMutex
}

func (p *Parser) ParseFile(name string) (*Document, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	f, err := p.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return p.Parse(f)
}

func (p *Parser) Parse(r io.Reader) (*Document, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	// pre parse
	r, err := p.PreParsers.PreParse(p, r)
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

	doc := p.newDoc()
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

func (p *Parser) ParseExecuteFragment(ctx context.Context, r io.Reader) (Nodes, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	nodes, err := p.ParseFragment(r)
	if err != nil {
		return nil, err
	}

	doc := p.newDoc()
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

	pages := ByType[*Page](doc.Nodes)

	if len(pages) > 0 {
		return pages[0].Nodes, nil
	}

	body, err := doc.Body()
	if err != nil {
		return nil, err
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
		return p.element(node, parent)
	case html.ErrorNode:
		return nil, fmt.Errorf(node.Data)
	case html.TextNode:
		return TextNode(node.Data), nil
	}

	return nil, fmt.Errorf("unknown node type %v", node.Data)
}

func (p *Parser) element(node *html.Node, parent Node) (Node, error) {
	el := &Element{
		Attributes: ConvertHTMLAttrs(node.Attr),
		HTMLNode:   node,
		Parent:     parent,
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
		return nil, err
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
