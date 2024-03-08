package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type PreParseError struct {
	Err       error
	Contents  []byte
	Filename  string
	Root      string
	PreParser any
}

func (e PreParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"contents":   string(e.Contents),
		"error":      errForJSON(e.Err),
		"filename":   e.Filename,
		"pre_parser": toType(e.PreParser),
		"root":       e.Root,
		"type":       toType(e),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (e PreParseError) Error() string {
	sb := strings.Builder{}

	var lines []string

	fp := filepath.Join(e.Root, e.Filename)
	if len(fp) > 0 {
		lines = append(lines, fmt.Sprintf("filepath: %s", fp))
	}

	if e.PreParser != nil {
		lines = append(lines, fmt.Sprintf("pre_parser: %s", toType(e.PreParser)))
	}

	if e.Err != nil {
		lines = append(lines, fmt.Sprintf("error: %s", e.Err))
	}

	sb.WriteString(strings.Join(lines, "\n"))

	s := sb.String()

	return strings.TrimSpace(s)
}

func (e PreParseError) Unwrap() error {
	if _, ok := e.Err.(unwrapper); ok {
		return errors.Unwrap(e.Err)
	}

	return e.Err
}

func (e PreParseError) As(target any) bool {
	ex, ok := target.(*PreParseError)
	if !ok {
		return false
	}

	(*ex) = e
	return true
}

func (e PreParseError) Is(target error) bool {
	if _, ok := target.(PreParseError); ok {
		return true
	}

	return false
}
