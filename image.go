package hype

import (
	"context"
	"encoding/json"
	"io/fs"
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

	m["type"] = "hype.Image"

	return json.Marshal(m)
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
