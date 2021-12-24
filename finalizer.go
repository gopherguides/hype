package hype

// Finalizer is a function that is called when the parser is finished.
type Finalizer interface {
	Finalize(p *Parser) error
}
