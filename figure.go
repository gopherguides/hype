package hype

import (
	"encoding/json"
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
	// bb.WriteString(f.StartTag())
	fmt.Fprintf(bb, "<a id=%q></a>\n", strings.TrimPrefix(f.Link(), "#"))
	fmt.Fprintln(bb)
	bb.WriteString(f.Nodes.MD())
	// bb.WriteString(f.EndTag())

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

func (f *Figure) String() string {
	if f == nil || f.Element == nil {
		return "<figure></figure>"
	}

	bb := &strings.Builder{}
	bb.WriteString(f.StartTag())
	bb.WriteString("\n")
	bb.WriteString(f.Nodes.String())
	bb.WriteString("\n")
	bb.WriteString(f.EndTag())
	return bb.String()
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

	if len(el.Nodes) == 0 {
		return f, nil
	}

	// Use the nodes directly to preserve their types (e.g., custom nodes like <godoc>)
	f.Nodes = el.Nodes
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
