package hype

import (
	"io/fs"

	"golang.org/x/net/html"
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

func (i Image) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator("img", "image"))
	return i.Node.Validate(html.ElementNode, checks...)
}

func (i Image) ValidateFS(cab fs.FS, checks ...ValidatorFn) error {
	checks = append(checks, SourceValidator(cab, &i))

	return i.Validate(checks...)
}

func NewImage(cab fs.FS, node *Node) (*Image, error) {

	i := &Image{
		Node: node,
	}

	return i, i.ValidateFS(cab)
}

func (p *Parser) NewImage(node *Node) (*Image, error) {
	return NewImage(p.FS, node)
}
