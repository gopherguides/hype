package hype

import (
	"fmt"
	"io/fs"

	"golang.org/x/net/html/atom"
)

var _ SetSourceable = &Image{}

type Image struct {
	*Node
}

func (c *Image) Source() (Source, bool) {
	c.RLock()
	defer c.RUnlock()
	return SrcAttr(c.attrs)
}

func (c *Image) SetSource(s string) {
	c.Lock()
	defer c.Unlock()
	c.attrs["src"] = s
}

func (i Image) String() string {
	if i.Node == nil {
		return "<img />"
	}
	return i.InlineTag()
}

func (p *Parser) NewImage(node *Node) (*Image, error) {
	return NewImage(p.FS, node)
}

func NewImage(cab fs.StatFS, node *Node) (*Image, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("image node can not be nil")
	}

	if !IsAtom(node, atom.Img) {
		return nil, fmt.Errorf("node is not an image %q", node.DataAtom)
	}

	i := &Image{
		Node: node,
	}

	source, ok := i.Source()
	if !ok {
		return nil, fmt.Errorf("image node has no src attribute")
	}

	if _, err := source.StatFile(cab); err != nil {
		return nil, err
	}

	return i, nil
}
