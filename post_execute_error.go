package hype

import (
	"encoding/json"
	"errors"
)

type PostExecuteError struct {
	Err          error
	Document     *Document
	Filename     string
	OrigErr      error
	Root         string
	PostExecuter any
}

func (pee PostExecuteError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"document":      pee.Document,
		"error":         errForJSON(pee.Err),
		"filename":      pee.Filename,
		"origal_error":  errForJSON(pee.OrigErr),
		"post_executer": toType(pee.PostExecuter),
		"root":          pee.Root,
		"type":          toType(pee),
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
