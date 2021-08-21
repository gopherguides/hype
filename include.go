package hype

import (
	"fmt"
	"path/filepath"
)

type Include struct {
	*Node
}

func (i Include) String() string {
	kids := i.Children
	if len(kids) > 0 {
		return i.Children.String()
	}

	if len(i.Data) > 0 {
		return i.Data
	}

	return fmt.Sprintf("<include %s />", i.Attrs())

}

func (p *Parser) NewInclude(node *Node) (*Include, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("include node can not be nil")
	}

	if node.Data != "include" {
		return nil, fmt.Errorf("node is not an include %q", node.Data)
	}

	i := &Include{
		Node: node,
	}

	src, err := i.Get("src")
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(src)

	switch ext {
	case ".html", ".md":
		// let these fall through as we'll handle them properly below
	default:
		b, err := p.ReadFile(src)
		if err != nil {
			return nil, err
		}
		i.Data = string(b)
		return i, nil
	}

	base := filepath.Base(src)
	dir := filepath.Dir(src)

	p2, err := p.SubParser(dir)
	if err != nil {
		return nil, err
	}

	doc, err := p2.ParseFile(base)
	if err != nil {
		return nil, err
	}

	body, err := doc.Body()
	if err != nil {
		return nil, err
	}

	i.Children = body.Children

	return i, nil
}
