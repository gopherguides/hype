package hype

import (
	"fmt"
	"path"
	"strings"

	"github.com/gopherguides/hype/atomx"
)

func NewGoDocLinkNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	ats, err := el.Attrs().Clone()
	if err != nil {
		return nil, err
	}

	af := NewEl(atomx.A, el)
	af.Attributes = ats

	href := el.Nodes.String()
	af.Set("href", fmt.Sprintf("https://pkg.go.dev/%s", href))
	af.Set("for", href)

	dt := strings.ReplaceAll(href, "#", ".")

	if _, ok := af.Get("full-pkg"); !ok {
		spl := strings.Split(href, "#")
		base := path.Base(spl[0])
		if len(spl) > 1 {
			x := spl[1]
			dt = fmt.Sprintf("%s.%s", base, x)

			if base == "builtin" {
				dt = x
			}
		}
	}

	cel := NewEl(atomx.Code, el)
	cel.Nodes = append(cel.Nodes, Text(dt))
	af.Nodes = append(af.Nodes, cel)

	l, err := NewLink(af)
	if err != nil {
		return nil, err
	}

	return Nodes{l}, nil
}
