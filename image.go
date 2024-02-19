package hype

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
)

type Image struct {
	*Element
}

func (i *Image) MarshalJSON() ([]byte, error) {
	if i == nil {
		return nil, ErrIsNil("image")
	}

	m, err := i.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", i)

	return json.MarshalIndent(m, "", "  ")
}

func (i *Image) MD() string {
	if i == nil {
		return ""
	}

	return i.Element.String()
}

func (i *Image) Execute(ctx context.Context, doc *Document) error {
	if i == nil {
		return ErrIsNil("image")
	}

	src, err := i.ValidAttr("src")
	if err != nil {
		return err
	}

	// Check if the file exists (but only if it's local)
	if strings.HasPrefix(src, "http") {
		return nil
	}
	if _, err := fs.Stat(doc.FS, src); err != nil {
		return err
	}

	return nil
}

func NewImage(el *Element) (*Image, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	i := &Image{
		Element: el,
	}

	if _, err := i.ValidAttr("src"); err != nil {
		return nil, err
	}

	return i, nil

}

func NewImageNodes(p *Parser, el *Element) (Nodes, error) {
	i, err := NewImage(el)
	if err != nil {
		return nil, err
	}

	return Nodes{i}, nil
}
