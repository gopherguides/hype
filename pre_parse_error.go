package hype

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PreParseError struct {
	Err       error
	Contents  []byte
	Filename  string
	Root      string
	PreParser any
}

func (e PreParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"contents":   string(e.Contents),
		"error":      errForJSON(e.Err),
		"filename":   e.Filename,
		"pre_parser": fmt.Sprintf("%T", e.PreParser),
		"root":       e.Root,
		"type":       fmt.Sprintf("%T", e),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (e PreParseError) Error() string {
	return toError(e)
}

func (e PreParseError) Unwrap() error {
	if _, ok := e.Err.(unwrapper); ok {
		return errors.Unwrap(e.Err)
	}

	return e.Err
}

func (e PreParseError) As(target any) bool {
	ex, ok := target.(*PreParseError)
	if !ok {
		return false
	}

	(*ex) = e
	return true
}

func (e PreParseError) Is(target error) bool {
	if _, ok := target.(PreParseError); ok {
		return true
	}

	return false
}
