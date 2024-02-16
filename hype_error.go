package hype

import (
	"encoding/json"
)

type HypeError struct {
	Err      error  `json:"err,omitempty"`
	Filename string `json:"filename,omitempty"`
	Root     string `json:"root,omitempty"`
}

func (he HypeError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"err":      he.Err,
		"filename": he.Filename,
		"root":     he.Root,
		"type":     "HypeError",
	}

	if _, ok := he.Err.(json.Marshaler); !ok && he.Err != nil {
		mm["err"] = he.Err.Error()
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (he HypeError) Error() string {
	if he.Err == nil {
		return ""
	}

	b, _ := json.MarshalIndent(he, "", "  ")
	return string(b)
}

func (he HypeError) Unwrap() error {
	type Unwrapper interface {
		Unwrap() error
	}

	if _, ok := he.Err.(Unwrapper); ok {
		return he.Err
	}

	return nil
}

func (he HypeError) As(target any) bool {
	ex, ok := target.(*HypeError)
	if !ok {
		return false
	}

	(*ex) = he
	return true
}

func (he HypeError) Is(target error) bool {
	if _, ok := target.(HypeError); ok {
		return true
	}

	return false
}
