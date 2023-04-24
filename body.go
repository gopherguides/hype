package hype

import (
	"encoding/json"
	"fmt"
)

// Body is a container for all the elements in a document.
type Body struct {
	*Element
}

func (b *Body) MarshalJSON() ([]byte, error) {
	if b == nil {
		return nil, ErrIsNil("body")
	}

	b.RLock()
	defer b.RUnlock()

	m, err := b.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", b)

	return json.Marshal(m)
}

// AsPage returns the body as a Page.
func (b *Body) AsPage() *Page {
	return &Page{
		Element: b.Element,
	}
}

// NewBody creates a new Body.
func NewBody(el *Element) (*Body, error) {
	if el == nil {
		return nil, el.WrapErr(ErrIsNil("element"))
	}

	body := &Body{
		Element: el,
	}

	return body, nil
}

// NewBodyNodes implements the ParseElementFn type
func NewBodyNodes(p *Parser, el *Element) (Nodes, error) {
	body, err := NewBody(el)
	if err != nil {
		return nil, err
	}

	return Nodes{body}, nil
}
