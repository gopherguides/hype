package hype

import (
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/gopherguides/hype/atomx"
)

var goVersion = runtime.Version

func GoVersion() string {
	return goVersion()
}

func NewGolangs(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	ats := el.Attrs()

	if err := setGolangAtrs(ats); err != nil {
		return nil, err
	}

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

				if err := el.Set("hide-cmd", ""); err != nil {
					return nil, err
				}

				el.Delete("data-go-version")
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

func setGolangAtrs(ats *Attributes) error {
	if ats == nil {
		return nil
	}

	if err := ats.Set("data-go-version", GoVersion()); err != nil {
		return err
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
		if err := ats.Set("environ", ev); err != nil {
			return err
		}
	}

	return nil
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
