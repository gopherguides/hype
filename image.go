package hype

import (
	"context"
	"io/fs"
)

type Image struct {
	*Element
}

func (i *Image) Execute(ctx context.Context, doc *Document) error {
	if i == nil {
		return ErrIsNil("image")
	}

	src, ok := i.Get("src")
	if !ok {
		return ErrAttrNotFound("src")
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

	src, ok := i.Get("src")
	if !ok {
		return nil, ErrAttrNotFound("src")
	}

	if len(src) == 0 {
		return nil, ErrAttrEmpty("src")
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
