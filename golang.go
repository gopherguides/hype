package hype

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/gopherguides/hype/atomx"
)

func NewGolangs(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	ats := el.Attrs()

	setGolangAtrs(ats)

	var cmds []string

	nodes, err := NewAttrCode(p, el)
	if err != nil {
		return nil, err
	}

	for gk, gv := range goCmds {
		v, ok := ats.Get(gk)
		if !ok {
			continue
		}

		for _, v := range strings.Split(v, ",") {

			if gk == "sym" {
				gk = "doc"
				if err := el.Set("language", "go"); err != nil {
					return nil, err
				}
			}

			if len(gv) >= 0 {
				gk = fmt.Sprintf("%s %s", gk, gv)
			}

			gk = strings.TrimSpace(gk)

			c := fmt.Sprintf("go %s %s", gk, v)
			c = strings.TrimSpace(c)
			cmds = append(cmds, c)
		}
	}

	sort.Strings(cmds)

	for _, v := range cmds {
		el := NewEl(atomx.Cmd, el)
		el.Attributes, err = ats.Clone()
		if err != nil {
			return nil, err
		}

		if err := el.Set("exec", v); err != nil {
			return nil, err
		}

		cmd, err := NewCmd(el)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, cmd)
	}

	return nodes, nil
}

func NewGolangNodes(p *Parser, el *Element) (Nodes, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	if el == nil {
		return nil, ErrIsNil("element")
	}

	return NewGolangs(p, el)
}

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
	cel.Nodes = append(cel.Nodes, TextNode(dt))
	af.Nodes = append(af.Nodes, cel)

	l, err := NewLink(af)
	if err != nil {
		return nil, err
	}

	return Nodes{l}, nil
}

func setGolangAtrs(ats *Attributes) {
	if ats == nil {
		return
	}

	env := []string{}

	if e, ok := ats.Get("environ"); ok {
		e = strings.TrimSpace(e)
		env = append(env, e)
	}

	for _, k := range []string{"GOOS", "GOARCH"} {
		if v, ok := ats.Get(strings.ToLower(k)); ok {
			v = strings.TrimSpace(v)
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	ev := strings.Join(env, ",")
	if len(ev) > 0 {
		ats.Set("environ", ev)
	}

}

var goCmds = map[string]string{
	"bug":      "",
	"build":    "",
	"clean":    "",
	"doc":      "",
	"env":      "",
	"fix":      "",
	"fmt":      "",
	"generate": "",
	"get":      "",
	"help":     "",
	"install":  "",
	"list":     "",
	"mod":      "",
	"run":      "",
	"sym":      "-cmd -u -src -short",
	"test":     "",
	"tool":     "",
	"version":  "",
	"vet":      "",
}
