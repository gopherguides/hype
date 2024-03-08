package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type ParseError struct {
	Err      error
	Filename string
	Root     string
	Contents []byte
}

func (pe ParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"contents": string(pe.Contents),
		"error":    errForJSON(pe.Err),
		"filename": pe.Filename,
		"root":     pe.Root,
		"type":     toType(pe),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (pe ParseError) Error() string {
	sb := &strings.Builder{}

	var lines []string

	fp := filepath.Join(pe.Root, pe.Filename)
	if len(fp) > 0 {
		lines = append(lines, fmt.Sprintf("filepath: %s", fp))
	}

	if pe.Err != nil {
		lines = append(lines, fmt.Sprintf("error: %s", pe.Err))
	}

	sb.WriteString(strings.Join(lines, "\n"))

	s := sb.String()

	return strings.TrimSpace(s)
}

func (pe ParseError) Unwrap() error {
	if _, ok := pe.Err.(unwrapper); ok {
		return errors.Unwrap(pe.Err)
	}

	return pe.Err
}

func (pe ParseError) As(target any) bool {
	ex, ok := target.(*ParseError)
	if !ok {
		return errors.As(pe.Err, target)
	}

	(*ex) = pe
	return true
}

func (pe ParseError) Is(target error) bool {
	if _, ok := target.(ParseError); ok {
		return true
	}

	return errors.Is(pe.Err, target)
}
