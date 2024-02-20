package hype

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ExecuteError struct {
	Err      error
	Filename string
	Root     string
	Contents []byte
}

func (pe ExecuteError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"type":     fmt.Sprintf("%T", pe),
		"error":    errForJSON(pe.Err),
		"root":     pe.Root,
		"filename": pe.Filename,
		"contents": string(pe.Contents),
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
