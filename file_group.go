package hype

import (
	"bytes"
	"fmt"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type FileGroup struct {
	*Node
	name string
}

func (FileGroup) Atom() atom.Atom {
	return atomx.FileGroup
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
		DataValidator("filegroup"),
		AttrValidator(Attributes{
			"name": "*",
		},
		),
	)
	return fg.Node.Validate(html.ElementNode, checks...)
}

func (p *Parser) NewFileGroup(node *Node) (*FileGroup, error) {

	fg := &FileGroup{
		Node: node,
	}

	err := fg.Validate()

	if err != nil {
		return nil, err
	}

	fg.Node.DataAtom = fg.Atom()

	return fg, fg.Validate()
}
