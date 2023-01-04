package hype

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gobuffalo/flect"
)

type Figure struct {
	*Element

	Pos       int
	SectionID int

	style string
}

// Type returns type of the figure.
// ex: "figure", "table", "listing", ...
func (f *Figure) Type() string {
	if f == nil || len(f.style) == 0 {
		return "figure"
	}

	return f.style
}

func (f *Figure) Name() string {
	if f == nil {
		return ""
	}

	style := flect.Titleize(f.Type())

	return fmt.Sprintf("%s %d.%d", style, f.SectionID, f.Pos)
}

func (f *Figure) Link() string {
	if f == nil {
		return ""
	}

	id, _ := f.Get("id")
	return fmt.Sprintf("#%s", id)
}

func NewFigure(p *Parser, el *Element) (*Figure, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	if p == nil {
		return nil, ErrIsNil("parser")
	}

	f := &Figure{
		Element: el,
	}

	if _, err := f.ValidAttr("id"); err != nil {
		return nil, err
	}

	f.style = "figure"
	style, ok := f.Get("type")
	style = strings.TrimSpace(style)

	if ok && len(style) > 0 {
		f.style = style
	}

	body := f.Nodes.String()
	body = strings.TrimSpace(body)

	if len(body) == 0 {
		return f, nil
	}

	p2, err := p.Sub(".")
	if err != nil {
		return nil, f.WrapErr(err)
	}

	// p2.NodeParsers[atomx.P] = func(p *Parser, el *Element) (Nodes, error) {
	// 	nodes, err := NewParagraphNodes(p, el)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if len(nodes) == 0 {
	// 		fmt.Println("no p nodes")
	// 		return nil, nil
	// 	}

	// 	x := nodes.String()
	// 	x = strings.TrimSpace(x)
	// 	if len(x) == 0 {
	// 		fmt.Println("no p nodes")
	// 		return nil, nil
	// 	}

	// 	fmt.Printf("TODO >> figure.go:93 nodes.String() %[1]T %+[1]v\n", nodes.String())

	// 	return nodes, nil
	// }

	nodes, err := p2.ParseFragment(strings.NewReader(body))
	if err != nil {
		if !errors.Is(err, ErrNilFigure) {
			return nil, f.WrapErr(err)
		}
	}

	pages := ByType[*Page](nodes)
	if len(pages) == 0 {
		f.Nodes = nodes

		return f, nil
	}

	page := pages[0]

	f.Nodes = page.Nodes

	return f, nil
}

func NewFigureNodes(p *Parser, el *Element) (Nodes, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	f, err := NewFigure(p, el)
	if err != nil {
		return nil, err
	}

	f.Lock()
	f.SectionID = p.Section
	f.Unlock()

	return Nodes{f}, nil
}
