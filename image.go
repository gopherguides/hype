package hype

import (
	"fmt"

	"golang.org/x/net/html/atom"
)

type Image struct {
	*Node
}

func (i Image) String() string {
	return i.InlineTag()
}

func (p *Parser) NewImage(node *Node) (*Image, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("image node can not be nil")
	}

	if node.DataAtom != atom.Img {
		return nil, fmt.Errorf("node is not an image %q", node.DataAtom.String())
	}

	i := &Image{
		Node: node,
	}

	src, err := i.Get("src")
	if err != nil {
		return nil, err
	}

	if _, err := p.Stat(src); err != nil {
		return nil, err
	}

	return i, nil
}
