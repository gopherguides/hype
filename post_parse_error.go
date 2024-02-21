package hype

import (
	"encoding/json"
	"errors"
)

type PostParseError struct {
	Err        error
	Document   *Document
	Filename   string
	OrigErr    error
	Root       string
	PostParser any
}

func (ppe PostParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"document":     ppe.Document,
		"error":        toError(ppe.Err),
		"filename":     ppe.Filename,
		"origal_error": toError(ppe.OrigErr),
		"post_parser":  toType(ppe.PostParser),
		"root":         ppe.Root,
		"type":         toType(ppe),
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
