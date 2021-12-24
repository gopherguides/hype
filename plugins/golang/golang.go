package golang

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/plugins/commander"
)

// Register registers all of this plugin's commands.
func Register(p *hype.Parser) {
	p.SetCustomTag(LINK, func(node *hype.Node) (hype.Tag, error) {
		return NewLink(node)
	})

	p.SetCustomTag(GO, func(node *hype.Node) (hype.Tag, error) {
		return NewGo(node)
	})
}

// NewGo returns a new commander.Cmd based on the given node.
// See commander.Cmd for more information.
//
// Example:
// 	"<go build="-o ./bin" timeout="10s" src="./cmd/foo"></go>"
// 	"<go test="-v -cover ./..." timeout="10s" src="./cmd/foo"></go>"
func NewGo(node *hype.Node) (hype.Tag, error) {
	if node == nil {
		return nil, fmt.Errorf("node is nil")
	}

	node.DataAtom = commander.CMD
	ats := node.Attrs()

	var env []string
	if e, ok := ats["environ"]; ok {
		e = strings.TrimSpace(e)
		env = append(env, e)
	}

	for _, k := range []string{"GOOS", "GOARCH"} {
		if v, ok := ats[strings.ToLower(k)]; ok {
			v = strings.TrimSpace(v)
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	ev := strings.Join(env, ",")
	if len(ev) > 0 {
		node.Set("environ", ev)
	}

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

	return commander.NewCmd(node)
}

var goCmds = map[string]hype.Attributes{
	"bug":   hype.Attributes{},
	"build": hype.Attributes{},
	"clean": hype.Attributes{},
	"doc": hype.Attributes{
		"hide-duration": "true",
	},
	"env": hype.Attributes{
		"hide-duration": "true",
	},
	"fix": hype.Attributes{},
	"fmt": hype.Attributes{
		"hide-duration": "true",
	},
	"generate": hype.Attributes{},
	"get":      hype.Attributes{},
	"install":  hype.Attributes{},
	"list":     hype.Attributes{},
	"mod":      hype.Attributes{},
	"run":      hype.Attributes{},
	"test": hype.Attributes{
		"hide-duration": "true",
	},
	"tool": hype.Attributes{},
	"version": hype.Attributes{
		"hide-data": "true",
	},
	"vet": hype.Attributes{},
	"help": hype.Attributes{
		"hide-duration": "true",
	},
}
