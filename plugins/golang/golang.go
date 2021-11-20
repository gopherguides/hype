package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/plugins/commander"
)

const cacheDir = ".hype/golang"

func CachePath() (string, error) {
	root, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fp := filepath.Join(root, cacheDir, runtime.Version())
	return fp, nil
}

func Register(p *hype.Parser, root string) {
	p.SetCustomTag(GODOC, func(node *hype.Node) (hype.Tag, error) {
		ats := node.Attrs()

		ex := fmt.Sprintf("go doc %s", ats["exec"])
		node.Set("exec", ex)
		node.Set("hide-duration", "true")

		return commander.NewCmd(node, root)
	})

	p.SetCustomTag(LINK, func(node *hype.Node) (hype.Tag, error) {
		return NewLink(node)
	})

	p.SetCustomTag(GO, func(node *hype.Node) (hype.Tag, error) {
		return NewGo(node, root)
	})
}

func NewGo(node *hype.Node, root string) (hype.Tag, error) {
	if node == nil {
		return nil, fmt.Errorf("node is nil")
	}

	node.DataAtom = commander.CMD
	ats := node.Attrs()

	var env []string
	if e, err := ats.Get("environ"); err == nil {
		env = append(env, e)
	}

	if goos, err := ats.Get("goos"); err == nil {
		if len(goos) > 0 {
			env = append(env, "GOOS="+goos)
		}
	}

	if goarch, err := ats.Get("goarch"); err == nil {
		if len(goarch) > 0 {
			env = append(env, "GOARCH="+goarch)
		}
	}

	node.Set("env", strings.Join(env, ","))

	for k, gats := range goCmds {
		ex, ok := ats[k]
		if !ok {
			continue
		}

		s := fmt.Sprintf("go %s %s", k, ex)
		for k, v := range gats {
			node.Set(k, v)
		}

		node.Set("exec", s)
		node.Delete(k)
	}

	return commander.NewCmd(node, root)
}

var goCmds = map[string]hype.Attributes{
	"bug":   nil,
	"build": nil,
	"clean": nil,
	"doc": hype.Attributes{
		"hide-duration": "true",
	},
	"env": hype.Attributes{
		"hide-duration": "true",
	},
	"fix": nil,
	"fmt": hype.Attributes{
		"hide-duration": "true",
	},
	"generate": nil,
	"get":      nil,
	"install":  nil,
	"list":     nil,
	"mod":      nil,
	"run":      nil,
	"test": hype.Attributes{
		"hide-duration": "true",
	},
	"tool": nil,
	"version": hype.Attributes{
		"hide-data": "true",
	},
	"vet": nil,
	"help": hype.Attributes{
		"hide-duration": "true",
	},
}
