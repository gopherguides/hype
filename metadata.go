package hype

import (
	"fmt"
	"strings"

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
