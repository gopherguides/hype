package hype

import (
	"bytes"
	"fmt"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

var _ Tag = &FileGroup{}
var _ Validatable = &FileGroup{}

// FileGroup represents a collection of files.
//
// HTML Attributes:
// 	name (required): The name of the file group.
type FileGroup struct {
	*Node
	name string
}

// Name returns the name of the file group.
func (fg *FileGroup) Name() string {
	fg.Lock()
	name := fg.name
	if len(name) == 0 {
		name = fg.attrs["name"]
		fg.name = name
	}
	fg.Unlock()
	return name
}

func (fg *FileGroup) String() string {
	fg.RLock()
	defer fg.RUnlock()

	bb := &bytes.Buffer{}

	fmt.Fprint(bb, fg.StartTag())
	fmt.Fprint(bb, fg.Children.String())
	fmt.Fprint(bb, fg.EndTag())

	return bb.String()
}

// Validate the file group
func (fg FileGroup) Validate(p *Parser, checks ...ValidatorFn) error {
	checks = append(checks,
		AtomValidator(atomx.Filegroup),
		AttrValidator(Attributes{
			"name": "*",
		},
		),
	)
	return fg.Node.Validate(p, html.ElementNode, checks...)
}

// NewFileGroup returns a new FileGroup from the given node.
func NewFileGroup(node *Node) (*FileGroup, error) {

	fg := &FileGroup{
		Node: node,
	}

	return fg, nil
}
