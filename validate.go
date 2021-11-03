package hype

import (
	"fmt"
	"io/fs"
	"net/http"

	"golang.org/x/net/html/atom"
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

func AtomValidator(a atom.Atom) ValidatorFn {
	return func(n *Node) error {

		if !IsAtom(n, a) {
			return fmt.Errorf("expected atom %q, got %q", a.String(), n.Atom().String())
		}

		return nil
	}
}

func DataValidator(data string) ValidatorFn {
	return func(n *Node) error {

		if n.Data != data {
			return fmt.Errorf("expected data %q, got %q", data, n.Data)
		}

		return nil
	}
}

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

func AttrValidator(query Attributes) ValidatorFn {
	return func(n *Node) error {
		if !n.Attrs().Matches(query) {
			return fmt.Errorf("%s: attributes did not match query: %v != %v", n.InlineTag(), n.Attrs(), query)
		}
		return nil
	}
}

func ChildrenValidators(tag Tag, checks ...ValidatorFn) []ValidatorFn {
	fn := func(n *Node) error {
		return n.Children.Validate(checks...)
	}

	chock := make([]ValidatorFn, len(checks))
	copy(chock, checks)

	chock = append(chock, fn)

	return chock

}
