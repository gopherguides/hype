package hype

import (
	"fmt"
	"path/filepath"

	"golang.org/x/net/html/atom"
)

type Include struct {
	*Node
}

func (i *Include) Atom() atom.Atom {
	return Include_Atom
}

func (c *Include) Source() (Source, bool) {
	c.RLock()
	defer c.RUnlock()
	return SrcAttr(c.attrs)
}

func (i Include) String() string {
	if i.Node == nil {
		return "<include />"
	}

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

	node.DataAtom = i.Atom()

	source, ok := i.Source()
	if !ok {
		return nil, fmt.Errorf("include node has no source")
	}

	ext := source.Ext()
	src := source.String()

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

	base := source.Base()
	dir := source.Dir()

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

	i.setSources(dir, body.Children)
	i.Children = body.Children

	return i, nil
}

func (i *Include) setSources(dir string, tags Tags) {
	for _, tag := range tags {
		i.setSources(dir, tag.GetChildren())

		st, ok := tag.(SetSourceable)
		if !ok {
			continue
		}

		source, ok := st.Source()
		if !ok {
			continue
		}

		xs := filepath.Join(dir, source.String())

		st.SetSource(xs)
	}
}
