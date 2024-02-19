package hype

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PostParseError struct {
	Err        error
	Filename   string
	OrigErr    error
	Root       string
	PostParser any
}

func (ppe PostParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"error":        toError(ppe.Err),
		"origal_error": toError(ppe.OrigErr),
		"filename":     ppe.Filename,
		"root":         ppe.Root,
		"post_parser":  fmt.Sprintf("%T", ppe.PostParser),
		"type":         fmt.Sprintf("%T", ppe),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (ppe PostParseError) Error() string {
	b, _ := ppe.MarshalJSON()
	return string(b)
}

func (ppe PostParseError) Unwrap() error {
	if _, ok := ppe.Err.(unwrapper); ok {
		return errors.Unwrap(ppe.Err)
	}

	return ppe.Err
}

func (ppe PostParseError) As(target any) bool {
	ex, ok := target.(*PostParseError)
	if !ok {
		return false
	}

	(*ex) = ppe
	return true
}

func (ppe PostParseError) Is(target error) bool {
	if _, ok := target.(PostParseError); ok {
		return true
	}

	return false
}
