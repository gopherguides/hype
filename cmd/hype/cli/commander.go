package cli

import (
	"github.com/gopherguides/hype"
	"github.com/markbates/cleo"
)

type Commander = cleo.Commander

type ParserCommander interface {
	cleo.Commander
	SetParser(p *hype.Parser) error
}
