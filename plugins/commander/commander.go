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

func Register(p *hype.Parser, root string) {

	fn := func(node *hype.Node) (hype.Tag, error) {
		return NewCmd(node, root)
	}

	p.SetCustomTag(CMD, fn)
	p.SetCustomTag(COMMAND, fn)

}
