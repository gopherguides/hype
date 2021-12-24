package hype

import (
	"fmt"
	"io/fs"
	"net/http"
)

type ValidatorFn func(n *Node) error

type Validatable interface {
	Validate(validators ...ValidatorFn) error
}

type ValidatableFS interface {
	ValidateFS(cab fs.FS, validators ...ValidatorFn) error
}

type ValidatableHTTP interface {
	ValidateHTTP(client *http.Client) error
}

// AtomValidator returns a validator that checks that the node has the given atoms.
func AtomValidator(atoms ...Atom) ValidatorFn {
	return func(n *Node) error {

		if !IsAtom(n, atoms...) {
			return fmt.Errorf("expected atom(s) %q, got %q", atoms, n.Atom())
		}

		return nil
	}
}

// SourceValidator returns a validator that checks that the node has a
// src attribute that points to a file in the given cab.
func SourceValidator(cab fs.FS, tag Tag) ValidatorFn {
	return func(n *Node) error {

		source, ok := TagSource(tag)
		if !ok {
			return fmt.Errorf("expected tag %v to have source", tag)
		}

		_, err := source.StatFile(cab)

		if err != nil {
			return err
		}

		return nil
	}
}

// AttrValidator returns a validator that checks that the node has the given attributes.
func AttrValidator(query Attributes) ValidatorFn {
	return func(n *Node) error {
		if !n.Attrs().Matches(query) {
			return fmt.Errorf("%s: attributes did not match query: %v != %v", n.InlineTag(), n.Attrs(), query)
		}
		return nil
	}
}

// ChildrenValidator returns validators that validate a tags children.
func ChildrenValidators(tag Tag, checks ...ValidatorFn) []ValidatorFn {
	fn := func(n *Node) error {
		return n.Children.Validate(checks...)
	}

	chock := make([]ValidatorFn, len(checks))
	copy(chock, checks)

	chock = append(chock, fn)

	return chock

}
