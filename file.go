package hype

import (
	"bytes"
	"fmt"
	"io/fs"

	"golang.org/x/net/html"
)

type File struct {
	*Node
}

func (c *File) Source() (Source, bool) {
	c.RLock()
	defer c.RUnlock()
	return SrcAttr(c.attrs)
}

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

func (f File) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AdamValidator(File_Adam))

	return f.Node.Validate(html.ElementNode, checks...)
}

func (f File) ValidateFS(cab fs.FS, checks ...ValidatorFn) error {
	checks = append(checks, SourceValidator(cab, &f))

	return f.Validate(checks...)
}

func NewFile(cab fs.FS, node *Node) (*File, error) {
	fg := &File{
		Node: node,
	}

	return fg, fg.ValidateFS(cab)
}

func (p *Parser) NewFile(node *Node) (*File, error) {
	return NewFile(p.FS, node)
}
