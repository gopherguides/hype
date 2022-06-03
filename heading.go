package hype

import (
	"fmt"
	"strconv"
	"strings"
)

type Heading struct {
	*Element
	level int
}

func (h Heading) Level() int {
	return h.level
}

func (h Heading) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if len(h.FileName) > 0 {
			fmt.Fprintf(f, "file://%s: ", h.FileName)
		}
		fmt.Fprintf(f, "%s\n", h.String())
	default:
		fmt.Fprintf(f, "%s", h.String())
	}
}

func NewHeading(el *Element) (*Heading, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	h := &Heading{
		Element: el,
	}

	l := strings.ToLower(el.Atom().String())
	l = strings.TrimPrefix(l, "h")

	i, err := strconv.Atoi(l)
	if err != nil {
		return nil, err
	}

	h.level = i

	return h, nil
}

func NewHeadingNodes(p *Parser, el *Element) (Nodes, error) {
	h, err := NewHeading(el)
	if err != nil {
		return nil, err
	}

	return Nodes{h}, nil
}
