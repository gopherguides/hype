package hype

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PostExecuteError struct {
	Err          error
	Filename     string
	OrigErr      error
	Root         string
	PostExecuter any
}

func (pee PostExecuteError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"error":         errForJSON(pee.Err),
		"origal_error":  errForJSON(pee.OrigErr),
		"filename":      pee.Filename,
		"root":          pee.Root,
		"post_executer": fmt.Sprintf("%T", pee.PostExecuter),
		"type":          fmt.Sprintf("%T", pee),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (e PostExecuteError) Error() string {
	return toError(e)
}

func (e PostExecuteError) Unwrap() error {
	if _, ok := e.Err.(unwrapper); ok {
		return errors.Unwrap(e.Err)
	}

	return e.Err
}

func (e PostExecuteError) As(target any) bool {
	ex, ok := target.(*PostExecuteError)
	if !ok {
		return false
	}

	(*ex) = e
	return true
}

func (e PostExecuteError) Is(target error) bool {
	if _, ok := target.(PostExecuteError); ok {
		return true
	}

	return false
}
