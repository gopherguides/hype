package hype

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

// File is a file node.
type File struct {
	*Node
}

// Source returns the source of the file.
func (c *File) Source() (Source, bool) {
	c.RLock()
	defer c.RUnlock()
	return SrcAttr(c.attrs)
}

// SetSource sets the source of the file.
func (c *File) SetSource(src string) {
	c.Lock()
	defer c.Unlock()

	c.attrs["src"] = src
}

func (f *File) String() string {
	if f.Node == nil {
		return "<file />"
	}

	bb := &bytes.Buffer{}

	fmt.Fprint(bb, f.StartTag())

	body := f.Children.String()

	fmt.Fprint(bb, body)

	fmt.Fprint(bb, f.EndTag())
	return bb.String()
}

// Validate the file
func (f File) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator(atomx.File))

	return f.Node.Validate(html.ElementNode, checks...)
}

// ValidateFS validates the file against the given filesystem.
func (f File) ValidateFS(cab fs.FS, checks ...ValidatorFn) error {
	checks = append(checks, SourceValidator(cab, &f))

	return f.Validate(checks...)
}

// NewFile returns a new File from the given node.
func NewFile(cab fs.FS, node *Node) (*File, error) {
	fg := &File{
		Node: node,
	}

	return fg, fg.ValidateFS(cab)
}
