package hype

type Finalizer interface {
	Finalize(p *Parser) error
}
