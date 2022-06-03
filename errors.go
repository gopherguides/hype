package hype

import (
	"fmt"
)

const (
	ErrNilFigure = ErrIsNil("figure")
)

type ErrIsNil string

func (e ErrIsNil) Error() string {
	return fmt.Sprintf("%s is nil", string(e))
}

func WrapNodeErr(n Node, err error) error {
	if err == nil {
		return nil
	}

	if h, ok := n.(Tag); ok {
		return fmt.Errorf("%T: %v: %w", h, h, err)
	}

	return fmt.Errorf("%T: %w", n, err)
}
