package hype

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"

	"golang.org/x/net/html/atom"
)

const (
	File_Atom atom.Atom = 1421757657
)

type File struct {
	*Node
}

func (c *File) Src() string {
	c.RLock()
	defer c.RUnlock()
	return c.attrs["src"]
}

func (c *File) SetSrc(src string) {
	c.Lock()
	defer c.Unlock()

	c.attrs["src"] = src
}

func (f File) String() string {
	bb := &bytes.Buffer{}

	fmt.Fprint(bb, f.StartTag())

	body := f.Children.String()

	fmt.Fprint(bb, body)
	fmt.Fprint(bb, f.EndTag())
	return bb.String()
}

func (p *Parser) NewFile(node *Node) (*File, error) {
	return NewFile(p.FS, node)
}

func NewFile(cab fs.StatFS, node *Node) (*File, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("file node can not be nil")
	}

	if node.Data != "file" {
		return nil, fmt.Errorf("node is not a file %q", node.Data)
	}

	node.DataAtom = File_Atom

	fg := &File{
		Node: node,
	}

	src, err := fg.Get("src")
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(src, "http") {
		return fg, nil
	}

	if _, err := cab.Stat(src); err != nil {
		return nil, err
	}

	return fg, nil
}
