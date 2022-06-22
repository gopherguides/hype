package cli

import (
	"fmt"

	"github.com/gopherguides/hype"
)

type Binding struct {
	*hype.Element

	Whole *Whole
}

func (b *Binding) String() string {
	if b == nil {
		return ""
	}

	return b.Children().String()
}

func NewBindingNode(el *hype.Element, whole *Whole) (hype.Node, error) {
	if el == nil {
		return nil, fmt.Errorf("element is nil")
	}

	if whole == nil {
		return nil, fmt.Errorf("whole is nil")
	}

	b := &Binding{
		Element: el,
		Whole:   whole,
	}

	if p, ok := b.Get("part"); ok {
		err := b.parseParts(p)
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	if p, ok := b.Get("whole"); ok {
		err := b.parseWholes(p)
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	return b, b.WrapErr(fmt.Errorf("binding is empty"))
}

func (b *Binding) parseWholes(p string) error {
	if b == nil {
		return b.WrapErr(hype.ErrIsNil("binding"))
	}

	whole := b.Whole

	if whole == nil {
		return b.WrapErr(hype.ErrIsNil("whole"))
	}

	switch p {
	case "title":
		b.Nodes = append(b.Nodes, hype.Text(whole.Name.Titleize().String()))
	case "":
		b.Nodes = append(b.Nodes, hype.Text(whole.Ident.String()))
	}

	return nil
}

func (b *Binding) parseParts(p string) error {

	if b == nil {
		return b.WrapErr(hype.ErrIsNil("binding"))
	}

	whole := b.Whole

	if whole == nil {
		return b.WrapErr(hype.ErrIsNil("whole"))
	}

	if len(p) == 0 {
		b.Nodes = append(b.Nodes, hype.Text(whole.PartIdent.String()))
		return nil
	}

	part, ok := whole.Parts[p]
	if !ok {
		return b.WrapErr(fmt.Errorf("part %q not found", p))
	}
	b.Nodes = append(b.Nodes, part)

	return nil
}

func NewBindingNodes(whole *Whole) hype.ParseElementFn {
	return func(p *hype.Parser, el *hype.Element) (hype.Nodes, error) {
		b, err := NewBindingNode(el, whole)
		if err != nil {
			return nil, err
		}

		return hype.Nodes{b}, nil
	}
}
