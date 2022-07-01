package cli

import (
	"fmt"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/binding"
)

type Binding struct {
	*hype.Element

	Whole *binding.Whole
}

func (b *Binding) String() string {
	if b == nil {
		return ""
	}

	return b.Children().String()
}

func NewBindingNode(p *hype.Parser, el *hype.Element, whole *binding.Whole) (hype.Node, error) {
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

	if x, ok := b.Get("part"); ok {
		err := b.parseParts(p, x)
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

func (b *Binding) parseWholes(arg string) error {
	if b == nil {
		return b.WrapErr(hype.ErrIsNil("binding"))
	}

	whole := b.Whole

	if whole == nil {
		return b.WrapErr(hype.ErrIsNil("whole"))
	}

	switch arg {
	case "title":
		b.Nodes = append(b.Nodes, hype.Text(whole.Name.Titleize().String()))
	case "":
		b.Nodes = append(b.Nodes, hype.Text(whole.Ident.String()))
	}

	return nil
}

func (b *Binding) parseParts(p *hype.Parser, key string) error {

	if b == nil {
		return b.WrapErr(hype.ErrIsNil("binding"))
	}

	whole := b.Whole

	if whole == nil {
		return b.WrapErr(hype.ErrIsNil("whole"))
	}

	if len(key) == 0 {
		b.Nodes = append(b.Nodes, hype.Text(whole.PartIdent.String()))
		return nil
	}

	part, ok := whole.Parts[key]
	if !ok {
		return b.WrapErr(fmt.Errorf("part %q not found", key))
	}

	p, err := p.Sub(part.Path)
	if err != nil {
		return b.WrapErr(err)
	}

	doc, err := p.ParseFile("module.md")
	if err != nil {
		return b.WrapErr(err)
	}

	part.Name = flect.New(doc.Title)

	whole.Parts[key] = part

	b.Nodes = append(b.Nodes, Part(part))

	return nil
}

func NewBindingNodes(whole *binding.Whole) hype.ParseElementFn {
	return func(p *hype.Parser, el *hype.Element) (hype.Nodes, error) {
		b, err := NewBindingNode(p, el, whole)
		if err != nil {
			return nil, err
		}

		return hype.Nodes{b}, nil
	}
}

type Part binding.Part

func (p Part) String() string {
	return fmt.Sprintf("\"%s %d: %s\"", p.Ident.Titleize(), p.Number, p.Name.Titleize())
}

func (p Part) Children() hype.Nodes {
	return nil
}
