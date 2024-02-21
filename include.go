package hype

import (
	"encoding/json"
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

func (inc *Include) MarshalJSON() ([]byte, error) {
	if inc == nil {
		return nil, ErrIsNil("include")
	}

	inc.RLock()
	defer inc.RUnlock()

	m, err := inc.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(inc)

	if inc.dir != "" {
		m["dir"] = inc.dir
	}

	return json.MarshalIndent(m, "", "  ")
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
	if inc == nil {
		return ""
	}
	return inc.Children().String()
}

func (inc *Include) MD() string {
	if inc == nil {
		return ""
	}

	return inc.Children().MD()
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

	fn := func(i int, fig *Figure) (string, error) {
		id, err := fig.ValidAttr("id")
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s.%d#%s", src, i, id), nil
	}

	if err := RestripeFigureIDs(inc.Nodes, fn); err != nil {
		return nil, ParseError{
			Err:      err,
			Filename: p.Filename,
			Root:     p.Root,
		}
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
		inc.Nodes.updateFileName(inc.dir)

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
