package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
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
	sb := &strings.Builder{}

	var lines []string

	fp := filepath.Join(e.Root, e.Filename)
	if len(fp) > 0 {
		lines = append(lines, fmt.Sprintf("filepath: %s", fp))
	}

	if e.Document != nil && len(e.Document.Title) > 0 {
		lines = append(lines, fmt.Sprintf("document: %s", e.Document.Title))
	}

	if e.Err != nil {
		lines = append(lines, fmt.Sprintf("post execute error: %s", e.Err))
	}

	if e.OrigErr != nil {
		lines = append(lines, fmt.Sprintf("original error: %s", e.OrigErr))
	}

	sb.WriteString(strings.Join(lines, "\n"))

	s := sb.String()

	return strings.TrimSpace(s)
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
