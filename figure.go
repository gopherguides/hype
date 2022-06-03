package hype

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gobuffalo/flect"
)

type Figure struct {
	*Element

	Pos       int
	SectionID int

	style string
	once  sync.Once
}

func (f *Figure) Name() string {
	if f == nil {
		return ""
	}
	style := f.style
	if len(style) == 0 {
		style = "figure"
	}

	// f.RLock()
	// defer f.RUnlock()
	style = flect.Titleize(style)
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
		return nil, err
	}

	nodes, err := p2.ParseFragment(strings.NewReader(body))
	if err != nil {
		if !errors.Is(err, ErrNilFigure) {
			return nil, err
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
	f.SectionID = p.Section

	return Nodes{f}, nil
}

type IDGenerator func(fig *Figure) (string, error)

func RestripeFigureIDs(nodes Nodes, fn IDGenerator) error {
	if fn == nil {
		return ErrIsNil("IDGenerator")
	}

	figs := ByType[*Figure](nodes)

	for _, fig := range figs {

		fid, err := fig.ValidAttr("id")
		if err != nil {
			return err
		}

		uid, err := fn(fig)
		if err != nil {
			return err
		}

		if err := fig.Set("id", uid); err != nil {
			return err
		}

		refs := ByType[*Ref](nodes)
		for _, ref := range refs {
			rid, err := ref.ValidAttr("id")
			if err != nil {
				return err
			}

			if rid != fid {
				continue
			}

			ref.Figure = fig
			if err := ref.Set("id", uid); err != nil {
				return err
			}

		}

	}

	return nil
}
