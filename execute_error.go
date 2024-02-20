package hype

import (
	"encoding/json"
	"errors"
)

type ExecuteError struct {
	Err      error
	Contents []byte
	Document *Document
	Filename string
	Root     string
}

func (pe ExecuteError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"contents": string(pe.Contents),
		"document": pe.Document,
		"error":    errForJSON(pe.Err),
		"filename": pe.Filename,
		"root":     pe.Root,
		"type":     toType(pe),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (pe ExecuteError) Error() string {
	return toError(pe)
}

func (pe ExecuteError) Unwrap() error {
	if _, ok := pe.Err.(unwrapper); ok {
		return errors.Unwrap(pe.Err)
	}

	return pe.Err
}

func (pe ExecuteError) As(target any) bool {
	ex, ok := target.(*ExecuteError)
	if !ok {
		return errors.As(pe.Err, target)
	}

	(*ex) = pe
	return true
}

func (pe ExecuteError) Is(target error) bool {
	if _, ok := target.(ExecuteError); ok {
		return true
	}

	return errors.Is(pe.Err, target)
}
