package hype

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Heading struct {
	*Element
	level int
}

func (h Heading) MarshalJSON() ([]byte, error) {
	h.RLock()
	defer h.RUnlock()

	m, err := h.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", h)
	m["level"] = h.level

	return json.MarshalIndent(m, "", "  ")
}

func (h Heading) MD() string {
	x := strings.Repeat("#", h.level)

	return fmt.Sprintf("%s %s", x, h.Children().MD())
}

func (h Heading) Level() int {
	return h.level
}

func (h Heading) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if len(h.Filename) > 0 {
			fmt.Fprintf(f, "file://%s: ", h.Filename)
		}
		fmt.Fprintf(f, "%s", h.String())
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
