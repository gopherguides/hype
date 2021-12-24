package hype

import (
	"io/fs"

	"golang.org/x/net/html"
)

var _ SetSourceable = &Image{}

// Image represents an HTML image.
type Image struct {
	*Node
}

// Source returns the source of the image.
func (c *Image) Source() (Source, bool) {
	c.RLock()
	defer c.RUnlock()
	return SrcAttr(c.attrs)
}

// SetSource sets the source of the image.
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

// Validate the image
func (i Image) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator("img", "image"))
	return i.Node.Validate(html.ElementNode, checks...)
}

// ValidateFS validates the image against the given filesystem.
func (i Image) ValidateFS(cab fs.FS, checks ...ValidatorFn) error {
	checks = append(checks, SourceValidator(cab, &i))

	return i.Validate(checks...)
}

// NewImage returns a new Image from the given node.
func NewImage(cab fs.FS, node *Node) (*Image, error) {

	i := &Image{
		Node: node,
	}

	return i, i.ValidateFS(cab)
}
