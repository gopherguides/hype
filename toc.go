package hype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type ToC struct {
	*Element
}

func (toc *ToC) MarshalJSON() ([]byte, error) {
	if toc == nil {
		return nil, ErrIsNil("toc")
	}

	m, err := toc.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", toc)

	return json.MarshalIndent(m, "", "  ")
}

func (toc *ToC) PostExecute(ctx context.Context, doc *Document, err error) error {
	if err != nil {
		return nil
	}

	if err := toc.validate(); err != nil {
		return err
	}

	if doc == nil {
		return toc.WrapErr(fmt.Errorf("document is nil"))
	}

	nodes, err := GenerateToC(doc.Parser, doc.Children())
	if err != nil {
		return toc.WrapErr(err)
	}

	toc.Nodes = nodes

	headings := ByType[*Heading](doc.Children())
	for i, h := range headings {
		link := Text(fmt.Sprintf("<a id=\"heading-%d\"></a>%s", i, h.Children().String()))
		h.Nodes = Nodes{link}
	}

	return nil
}

func NewToCNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, fmt.Errorf("el is nil")
	}

	toc := &ToC{
		Element: el,
	}

	return Nodes{toc}, nil
}

func GenerateToC(p *Parser, nodes Nodes) (Nodes, error) {
	headings := ByType[*Heading](nodes)

	levels := map[int]int{}

	bb := &bytes.Buffer{}

	for i, h := range headings {
		t := h.Children().String()

		for i := 1; i < h.Level(); i++ {
			fmt.Fprint(bb, "\t")
		}

		switch h.Level() {
		case 1:
			x := levels[1]
			x++

			levels = map[int]int{
				1: x,
			}
		default:
			levels[h.Level()]++
		}

		dots := []string{fmt.Sprint(p.Section)}
		for i := 1; i <= h.level; i++ {
			dots = append(dots, fmt.Sprint(levels[i]))
		}

		lvl := strings.Join(dots, ".")
		lvl = fmt.Sprintf("<level>%s</level>", lvl)

		fmt.Fprintf(bb, "- <a href=\"#heading-%d\">%s %s</a>\n", i, lvl, t)
	}

	frag, err := p.ParseFragment(bb)
	if err != nil {
		return nil, err
	}

	for _, n := range frag.Children() {
		switch n.(type) {
		case Text:
		default:
			return Nodes{n}, nil
		}
	}

	return nil, fmt.Errorf("unable to parse toc")
}

func (toc *ToC) validate() error {
	if toc == nil {
		return toc.WrapErr(fmt.Errorf("toc is nil"))
	}

	if toc.Element == nil {
		return toc.WrapErr(fmt.Errorf("toc.Element is nil"))
	}

	return nil
}
