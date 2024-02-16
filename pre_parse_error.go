package hype

import (
	"encoding/json"
	"fmt"
)

type PreParseError struct {
	Err       error  `json:"error,omitempty"`
	Filename  string `json:"filename,omitempty"`
	Root      string `json:"root,omitempty"`
	PreParser PreParser
}

func (e PreParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"error":     e.Err,
		"filename":  e.Filename,
		"root":      e.Root,
		"preparser": fmt.Sprintf("%T", e.PreParser),
		"type":      fmt.Sprintf("%T", e),
	}

	if _, ok := e.Err.(json.Marshaler); !ok && e.Err != nil {
		mm["error"] = e.Err.Error()
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (e PreParseError) Error() string {
	b, err := e.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("error marshalling to json: %s", err)
	}
	return string(b)
}

func (e PreParseError) Unwrap() error {
	type Unwrapper interface {
		Unwrap() error
	}

	if _, ok := e.Err.(Unwrapper); ok {
		return e.Err
	}

	return nil
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
