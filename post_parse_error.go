package hype

import (
	"encoding/json"
	"fmt"
)

type PostParseError struct {
	Err        error      `json:"err,omitempty"`
	Filename   string     `json:"filename,omitempty"`
	OrigErr    error      `json:"orig_err,omitempty"`
	Root       string     `json:"root,omitempty"`
	PostParser PostParser `json:"-,omitempty"`
}

func (ppe PostParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"error":      ppe.Err,
		"filename":   ppe.Filename,
		"root":       ppe.Root,
		"postparser": fmt.Sprintf("%T", ppe.PostParser),
		"type":       fmt.Sprintf("%T", ppe),
	}

	if _, ok := ppe.Err.(json.Marshaler); !ok && ppe.Err != nil {
		mm["error"] = ppe.Err.Error()
	}

	if ppe.OrigErr != nil {
		mm["orig_error"] = ppe.OrigErr.Error()
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (ppe PostParseError) Error() string {
	b, _ := ppe.MarshalJSON()
	return string(b)
}

func (ppe PostParseError) Unwrap() error {
	type Unwrapper interface {
		Unwrap() error
	}

	if _, ok := ppe.Err.(Unwrapper); ok {
		return ppe.Err
	}

	return nil
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
