package hype

import (
	"bytes"
	"fmt"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

type FileGroup struct {
	*Node
	name string
}

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

func (fg FileGroup) Validate(checks ...ValidatorFn) error {
	checks = append(checks,
		AtomValidator(atomx.Filegroup),
		AttrValidator(Attributes{
			"name": "*",
		},
		),
	)
	return fg.Node.Validate(html.ElementNode, checks...)
}

func NewFileGroup(node *Node) (*FileGroup, error) {

	fg := &FileGroup{
		Node: node,
	}

	return fg, fg.Validate()
}
