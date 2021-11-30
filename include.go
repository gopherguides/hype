package hype

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

type Include struct {
	*Node
	Data string
}

func (c *Include) Source() (Source, bool) {
	c.RLock()
	defer c.RUnlock()
	return SrcAttr(c.attrs)
}

func (i Include) String() string {
	if i.Node == nil {
		return "<include />"
	}

	kids := i.Children
	if len(kids) > 0 {
		return i.Children.String()
	}

	return fmt.Sprintf("<include %s />", i.Attrs())
}

func (i Include) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator(atomx.Include))
	return i.Node.Validate(html.ElementNode, checks...)
}

func (i Include) ValidateFS(fs fs.FS, checks ...ValidatorFn) error {
	checks = append(checks, SourceValidator(fs, &i))
	return i.Validate(checks...)
}

func NewInclude(node *Node, p *Parser) (*Include, error) {

	i := &Include{
		Node: node,
	}

	if err := i.ValidateFS(p.FS); err != nil {
		return nil, err
	}

	source, ok := i.Source()
	if !ok {
		return nil, fmt.Errorf("include node has no source")
	}

	ext := source.Ext()
	src := source.String()

	switch ext {
	case ".html", ".md":
		// let these fall through as we'll handle them properly below
	default:
		b, err := p.ReadFile(src)
		if err != nil {
			return nil, err
		}
		i.Data = string(b)
		return i, nil
	}

	base := source.Base()
	dir := source.Dir()

	p2, err := p.SubParser(dir)
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

	i.setSources(dir, body.Children)
	i.Children = body.Children

	return i, nil
}

func (i *Include) setSources(dir string, tags Tags) {
	for _, tag := range tags {
		i.setSources(dir, tag.GetChildren())

		srcs := tag.GetChildren().ByAttrs(Attributes{
			"src": "*",
		})

		for _, src := range srcs {
			st, ok := src.(SetSourceable)
			if !ok {
				continue
			}

			source, ok := st.Source()
			if !ok {
				continue
			}

			ss := source.String()
			if strings.HasPrefix(ss, dir) {
				continue
			}

			xs := filepath.Join(dir, source.String())

			st.SetSource(xs)
		}
	}
}
