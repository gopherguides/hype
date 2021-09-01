package hype

import (
	"fmt"
	"io/fs"
	"strings"

	"golang.org/x/net/html/atom"
)

type Image struct {
	*Node
}

func (c *Image) Src() string {
	c.RLock()
	defer c.RUnlock()
	return c.attrs["src"]
}

func (i Image) String() string {
	return i.InlineTag()
}

func (p *Parser) NewImage(node *Node) (*Image, error) {
	return NewImage(p.FS, node)
}

func NewImage(cab fs.StatFS, node *Node) (*Image, error) {
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

	if strings.HasPrefix(src, "http") {
		return i, nil
	}

	if _, err := cab.Stat(src); err != nil {
		return nil, err
	}

	return i, nil
}
