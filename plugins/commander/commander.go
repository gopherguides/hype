package commander

import (
	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/atomx"
)

const (
	CMD     atomx.Atom = "cmd"
	COMMAND atomx.Atom = "command"
)

type Attributes = hype.Attributes
type Data map[string]string

// Register registers all of the atoms and tags
// that this plugin provides.
func Register(p *hype.Parser) {

	fn := func(node *hype.Node) (hype.Tag, error) {
		return NewCmd(node)
	}

	p.SetCustomTag(CMD, fn)
	p.SetCustomTag(COMMAND, fn)

}

// cmds is a map of special data attributes for different
// exacutable commands.
var cmds = map[string]hype.Attributes{
	"tree": hype.Attributes{
		"hide-data": "true",
	},
	"cat": hype.Attributes{
		"hide-data": "true",
	},
}
