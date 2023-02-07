package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"github.com/markbates/syncx"
)

type Metadata struct {
	Element *Element
	syncx.Map[string, string]
}

func (md *Metadata) IsEmptyNode() bool {
	return false
}

func (md *Metadata) Children() Nodes {
	if md == nil || md.Element == nil {
		return nil
	}

	return md.Element.Children()
}

func (md *Metadata) PostParse(p *Parser, d *Document, err error) error {
	if md == nil {
		return fmt.Errorf("metadata is nil")
	}

	if d == nil {
		return fmt.Errorf("document is nil")
	}

	heads := ByAtom(d.Nodes, atomx.Head)
	if len(heads) == 0 {
		return nil
	}

	hd := heads[0]

	head, ok := hd.(*Element)
	if !ok {
		return fmt.Errorf("head is not an element: %T", hd)
	}

	md.Map.Range(func(key, value string) bool {
		el := NewEl(atomx.Meta, head)
		el.Set(key, value)
		head.Nodes = append(head.Nodes, el)
		return true
	})

	return nil
}

func NewMetadata(el *Element) (*Metadata, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	m := &Metadata{
		Element: el,
		Map:     syncx.Map[string, string]{},
	}

	body := m.Children().String()
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, el.WrapErr(fmt.Errorf("invalid metadata line: %s", line))
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if err := m.Set(key, value); err != nil {
			return nil, el.WrapErr(err)
		}
	}

	m.Element.Nodes = Nodes{}

	return m, nil
}

func NewMetadataNodes(p *Parser, el *Element) (Nodes, error) {
	md, err := NewMetadata(el)
	if err != nil {
		return nil, err
	}

	return Nodes{md}, nil
}
