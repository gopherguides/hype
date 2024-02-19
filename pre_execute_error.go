package hype

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PreExecuteError struct {
	Err         error
	Filename    string
	Root        string
	PreExecuter any
}

func (pee PreExecuteError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"error":        errForJSON(pee.Err),
		"filename":     pee.Filename,
		"root":         pee.Root,
		"pre_executer": fmt.Sprintf("%T", pee.PreExecuter),
		"type":         fmt.Sprintf("%T", pee),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (pee PreExecuteError) Error() string {
	return toError(pee)
}

func (pee PreExecuteError) Unwrap() error {
	if _, ok := pee.Err.(unwrapper); ok {
		return errors.Unwrap(pee.Err)
	}

	return pee.Err
}

func (pee PreExecuteError) As(target any) bool {
	ex, ok := target.(*PreExecuteError)
	if !ok {
		return false
	}

	(*ex) = pee
	return true
}

func (pee PreExecuteError) Is(target error) bool {
	if _, ok := target.(PreExecuteError); ok {
		return true
	}

	return errors.Is(pee.Err, target)
}