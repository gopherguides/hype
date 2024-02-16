package hype

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ExecuteError struct {
	HypeError
}

func (pe ExecuteError) MarshalJSON() ([]byte, error) {
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

func (pe ExecuteError) Error() string {
	if pe.Err == nil {
		return ""
	}

	b, _ := json.MarshalIndent(pe, "", "  ")
	return string(b)
}

func (pe ExecuteError) Unwrap() error {
	type Unwrapper interface {
		Unwrap() error
	}

	if _, ok := pe.Err.(Unwrapper); ok {
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
