package hype

import (
	"context"
	"sync"

	"github.com/gopherguides/hype/atomx"
)

type Ref struct {
	*Element
	*Figure
	once sync.Once
}

func (r *Ref) PostExecute(ctx context.Context, doc *Document, err error) error {
	if err != nil {
		return nil
	}

	if r == nil {
		return ErrIsNil("ref")
	}

	if r.Element == nil {
		return ErrIsNil("element")
	}

	if r.Figure == nil {
		return ErrIsNil("figure")
	}

	href := NewEl(atomx.A, r)
	if err := href.Set("href", r.Link()); err != nil {
		return err
	}

	href.Nodes = append(href.Nodes, TextNode(r.Figure.Name()))
	r.Nodes = append(r.Nodes, href)

	return nil
}

func NewRef(el *Element) (*Ref, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	r := &Ref{
		Element: el,
	}

	id, ok := el.Get("id")
	if !ok {
		return nil, ErrAttrNotFound("id")
	}

	if len(id) == 0 {
		return nil, ErrAttrEmpty("id")
	}

	return r, nil
}

func NewRefNodes(p *Parser, el *Element) (Nodes, error) {
	r, err := NewRef(el)
	if err != nil {
		return nil, err
	}

	return Nodes{r}, nil
}
