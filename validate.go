package hype

import (
	"fmt"
	"io/fs"
	"net/http"
)

type ValidatorFn func(p *Parser, n *Node) error

type Validatable interface {
	Validate(p *Parser, validators ...ValidatorFn) error
}

type ValidatableFS interface {
	ValidateFS(p *Parser, cab fs.FS, validators ...ValidatorFn) error
}

type ValidatableHTTP interface {
	ValidateHTTP(client *http.Client) error
}

// AtomValidator returns a validator that checks that the node has the given atoms.
func AtomValidator(atoms ...Atom) ValidatorFn {
	return func(p *Parser, n *Node) error {

		if !IsAtom(n, atoms...) {
			return fmt.Errorf("expected atom(s) %q, got %q", atoms, n.Atom())
		}

		return nil
	}
}

// SourceValidator returns a validator that checks that the node has a
// src attribute that points to a file in the given cab.
func SourceValidator(cab fs.FS, tag Tag) ValidatorFn {
	return func(p *Parser, n *Node) error {

		source, ok := TagSource(tag)
		if !ok {
			return fmt.Errorf("expected tag %v to have source", tag)
		}

		if source.IsFile() {
			if _, err := source.StatFile(cab); err != nil {
				return err
			}
			return nil
		}

		if p == nil {
			return nil
		}

		client := p.Client
		if client == nil {
			return fmt.Errorf("no http client available")
		}

		_, err := source.StatHTTP(client)

		if err != nil {
			return err
		}

		return nil
	}
}

// AttrValidator returns a validator that checks that the node has the given attributes.
func AttrValidator(query Attributes) ValidatorFn {
	return func(p *Parser, n *Node) error {
		if !n.Attrs().Matches(query) {
			return fmt.Errorf("%s: attributes did not match query: %v != %v", n.InlineTag(), n.Attrs(), query)
		}
		return nil
	}
}

// ChildrenValidator returns validators that validate a tags children.
func ChildrenValidators(tag Tag, p *Parser, checks ...ValidatorFn) []ValidatorFn {

	chock := make([]ValidatorFn, len(checks))
	copy(chock, checks)

	fn := func(p *Parser, n *Node) error {
		return n.Children.Validate(p, checks...)
	}
	chock = append(chock, fn)

	return chock

}
