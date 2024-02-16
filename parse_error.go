package hype

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ParseError struct {
	Err      error  `json:"error,omitempty"`
	Filename string `json:"filename,omitempty"`
	Root     string `json:"root,omitempty"`
}

func (pe ParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"type":     fmt.Sprintf("%T", pe),
		"err":      pe.Err,
		"root":     pe.Root,
		"filename": pe.Filename,
	}

	if _, ok := pe.Err.(json.Marshaler); !ok && pe.Err != nil {
		mm["err"] = pe.Err.Error()
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (pe ParseError) Error() string {
	if pe.Err == nil {
		return ""
	}

	b, _ := json.MarshalIndent(pe, "", "  ")
	return string(b)
}

func (pe ParseError) Unwrap() error {
	type Unwrapper interface {
		Unwrap() error
	}

	if _, ok := pe.Err.(Unwrapper); ok {
		return errors.Unwrap(pe.Err)
	}

	return pe.Err
}

func (pe ParseError) As(target any) bool {
	ex, ok := target.(*ParseError)
	if !ok {
		return errors.As(pe.Err, target)
	}

	(*ex) = pe
	return true
}

func (pe ParseError) Is(target error) bool {
	if _, ok := target.(ParseError); ok {
		return true
	}

	return errors.Is(pe.Err, target)
}