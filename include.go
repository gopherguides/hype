package hype

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

type Include struct {
	*Element

	dir string
	pp  sync.Once
}

func (inc *Include) PostParse(p *Parser, d *Document, err error) error {
	if err != nil {
		return nil
	}

	if inc == nil {
		return ErrIsNil("include")
	}

	if err := inc.setSources(); err != nil {
		return err
	}

	return nil
}

func (inc *Include) String() string {
	return inc.Children().String()
}

func NewInclude(p *Parser, el *Element) (*Include, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	if el == nil {
		return nil, ErrIsNil("element")
	}

	if p.FS == nil {
		return nil, ErrIsNil("p.FS")
	}

	src, ok := el.Get("src")
	if !ok {
		return nil, fmt.Errorf("missing src attribute")
	}

	sdir := filepath.Dir(src)
	base := filepath.Base(src)

	p2, err := p.Sub(sdir)
	if err != nil {
		return nil, err
	}

	doc, err := p2.ParseFile(base)
	if err != nil {
		return nil, err
	}

	body, err := doc.Body()
	if err != nil {
		return nil, err
	}

	inc := &Include{
		Element: el,
		dir:     sdir,
	}

	inc.Nodes = body.Nodes

	fn := func(fig *Figure) (string, error) {
		id, ok := fig.Get("id")
		if !ok || len(id) == 0 {
			return "", ErrAttrEmpty("id")
		}

		return fmt.Sprintf("%s#%s", src, id), nil
	}

	if err := RestripeFigureIDs(inc.Nodes, fn); err != nil {
		return nil, err
	}

	return inc, nil
}

// NewIncludeNodes implements the ParseElementFn type
func NewIncludeNodes(p *Parser, el *Element) (Nodes, error) {
	inc, err := NewInclude(p, el)
	if err != nil {
		return nil, err
	}

	return Nodes{inc}, nil
}

func (inc *Include) setSources() error {
	if inc == nil {
		return nil
	}

	var err error
	inc.pp.Do(func() {
		kids := ByAttrs(inc.Children(), map[string]string{
			"src": "*",
		})

		for _, n := range kids {
			ats := n.Attrs()

			src, _ := ats.Get("src")

			if strings.HasPrefix(src, "http") || strings.HasPrefix(src, inc.dir) {
				continue
			}

			src = filepath.Join(inc.dir, src)

			if err = ats.Set("src", src); err != nil {
				return
			}
		}

	})

	return err
}
