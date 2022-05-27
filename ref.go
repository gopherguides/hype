package hype

import (
	"context"
	"fmt"
	"strings"
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
		return fmt.Errorf("%w: %s", ErrIsNil("figure"), r.StartTag())
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

	// get the id from the inner nodes
	id := el.Nodes.String()

	// use the id attr, if it exists
	if i, ok := el.Get("id"); ok {
		id = i
	}

	id = strings.TrimSpace(id)

	if len(id) == 0 {
		return nil, ErrAttrEmpty("id")
	}

	// set the id back on the element
	// for consistency
	if err := r.Set("id", id); err != nil {
		return nil, err
	}

	// clear out any existing inner nodes
	r.Nodes = Nodes{}

	return r, nil
}

func NewRefNodes(p *Parser, el *Element) (Nodes, error) {
	r, err := NewRef(el)
	if err != nil {
		return nil, err
	}

	return Nodes{r}, nil
}
