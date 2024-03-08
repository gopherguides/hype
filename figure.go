package hype

import (
	"encoding/json"
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

func (f *Figure) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, ErrIsNil("figure")
	}

	f.RLock()
	defer f.RUnlock()

	m, err := f.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(f)
	m["pos"] = f.Pos
	m["section_id"] = f.SectionID
	m["style"] = f.Style()

	return json.MarshalIndent(m, "", "  ")
}

func (f *Figure) MD() string {
	if f == nil {
		return ""
	}

	bb := &strings.Builder{}
	bb.WriteString(f.StartTag())
	fmt.Fprintln(bb)
	bb.WriteString(f.Nodes.String())
	bb.WriteString(f.EndTag())

	return bb.String()
}

// Style returns type of the figure.
// ex: "figure", "table", "listing", ...
func (f *Figure) Style() string {
	if f == nil || len(f.style) == 0 {
		return "figure"
	}

	return f.style
}

func (f *Figure) Name() string {
	if f == nil {
		return ""
	}

	style := flect.Titleize(f.Style())

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
