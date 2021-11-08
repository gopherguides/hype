package golang

import (
	"github.com/gopherguides/hype"
)

func Register(p *hype.Parser) {
	p.SetCustomTag(GODOC, func(node *hype.Node) (hype.Tag, error) {
		return NewGodoc(node)
	})

	p.SetCustomTag(LINK, func(node *hype.Node) (hype.Tag, error) {
		return NewLink(node)
	})
}
