package hype

import (
	"encoding/json"
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

func errForJSON(err error) any {
	if err == nil {
		return nil
	}

	if _, ok := err.(json.Marshaler); ok {
		return err
	}

	return err.Error()
}

func toError(err error) string {
	if err == nil {
		return ""
	}

	if _, ok := err.(json.Marshaler); ok {
		b, err := json.MarshalIndent(err, "", "  ")
		if err != nil {
			return "error marshalling to json: " + err.Error()
		}
		return string(b)
	}

	return err.Error()
}
