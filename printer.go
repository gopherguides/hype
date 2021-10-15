package hype

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type TransformerFn func(tag Tag) (Tag, error)

func defaultTransformer(tag Tag) (Tag, error) {
	return tag, nil
}

type Printer struct {
	*sync.RWMutex
	former TransformerFn
	output io.Writer
}

func (p *Printer) Transformer() TransformerFn {
	p.Lock()
	if p.former == nil {
		p.former = defaultTransformer
	}
	fn := p.former
	p.Unlock()

	return fn
}

func (p *Printer) Print(tags ...Tag) error {

	type taggable interface {
		EndTag() string
		StartTag() string
	}

	for _, tag := range tags {

		tag, err := p.Transform(tag)
		if err != nil {
			return err
		}

		tb, ok := tag.(taggable)
		if !ok {
			return fmt.Errorf("cant print tag %v", tag)
		}

		fmt.Fprint(p.Out(), tb.StartTag())
		if err := p.Print(tag.GetChildren()...); err != nil {
			return err
		}
		fmt.Fprint(p.Out(), tb.EndTag())

	}
	return nil
}

func (p *Printer) Transform(tag Tag) (Tag, error) {
	return p.Transformer()(tag)
}

func (p *Printer) SetTransformer(fn TransformerFn) {
	p.Lock()
	defer p.Unlock()

	p.former = fn
}

func (p *Printer) SetOutput(w io.Writer) {
	p.Lock()
	defer p.Unlock()
	p.output = w
}

func (p *Printer) Out() io.Writer {
	p.RLock()
	defer p.RUnlock()

	if p.output != nil {
		return p.output
	}

	return os.Stdout
}

func NewPrinter(w io.Writer) *Printer {
	p := &Printer{
		RWMutex: &sync.RWMutex{},
		former:  defaultTransformer,
		output:  w,
	}

	return p
}
