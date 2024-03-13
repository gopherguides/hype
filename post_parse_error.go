package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type PostParseError struct {
	Err        error
	Document   *Document
	Filename   string
	OrigErr    error
	Root       string
	PostParser any
}

func (ppe PostParseError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"document":     ppe.Document,
		"error":        toError(ppe.Err),
		"filename":     ppe.Filename,
		"origal_error": toError(ppe.OrigErr),
		"post_parser":  toType(ppe.PostParser),
		"root":         ppe.Root,
		"type":         toType(ppe),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (ppe PostParseError) Error() string {
	sb := &strings.Builder{}

	var lines []string

	fp := filepath.Join(ppe.Root, ppe.Filename)
	if len(fp) > 0 {
		lines = append(lines, fmt.Sprintf("filepath: %s", fp))
	}

	if ppe.Document != nil && len(ppe.Document.Title) > 0 {
		lines = append(lines, fmt.Sprintf("document: %s", ppe.Document.Title))
	}

	if ppe.Err != nil {
		lines = append(lines, fmt.Sprintf("post parse error: %s", ppe.Err))
	}

	if ppe.OrigErr != nil {
		lines = append(lines, fmt.Sprintf("original error: %s", ppe.OrigErr))
	}

	sb.WriteString(strings.Join(lines, "\n"))

	s := sb.String()

	return strings.TrimSpace(s)
}

func (ppe PostParseError) Unwrap() error {
	if _, ok := ppe.Err.(unwrapper); ok {
		return errors.Unwrap(ppe.Err)
	}

	return ppe.Err
}

func (ppe PostParseError) As(target any) bool {
	ex, ok := target.(*PostParseError)
	if !ok {
		return false
	}

	(*ex) = ppe
	return true
}

func (ppe PostParseError) Is(target error) bool {
	if _, ok := target.(PostParseError); ok {
		return true
	}

	return false
}
