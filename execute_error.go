package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
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
	sb := &strings.Builder{}

	var lines []string

	fp := filepath.Join(pe.Root, pe.Filename)
	if len(fp) > 0 {
		lines = append(lines, fmt.Sprintf("filepath: %s", fp))
	}

	if pe.Document != nil && len(pe.Document.Title) > 0 {
		lines = append(lines, fmt.Sprintf("document: %s", pe.Document.Title))
	}

	if pe.Err != nil {
		lines = append(lines, fmt.Sprintf("execute error: %s", pe.Err))
	}

	sb.WriteString(strings.Join(lines, "\n"))

	s := sb.String()

	return strings.TrimSpace(s)
}

func (pe ExecuteError) String() string {
	return pe.Error()
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
