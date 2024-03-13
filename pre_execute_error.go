package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type PreExecuteError struct {
	Err         error
	Document    *Document
	Filename    string
	Root        string
	PreExecuter any
}

func (pee PreExecuteError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"document":     pee.Document,
		"error":        errForJSON(pee.Err),
		"filename":     pee.Filename,
		"pre_executer": toType(pee.PreExecuter),
		"root":         pee.Root,
		"type":         toType(pee),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (pee PreExecuteError) Error() string {
	sb := &strings.Builder{}

	var lines []string

	fp := filepath.Join(pee.Root, pee.Filename)
	if len(fp) > 0 {
		lines = append(lines, fmt.Sprintf("filepath: %s", fp))
	}

	if pee.Document != nil && len(pee.Document.Title) > 0 {
		lines = append(lines, fmt.Sprintf("document: %s", pee.Document.Title))
	}

	if pee.Err != nil {
		lines = append(lines, fmt.Sprintf("pre execute error: %s", pee.Err))
	}

	sb.WriteString(strings.Join(lines, "\n"))

	s := sb.String()

	return strings.TrimSpace(s)
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
