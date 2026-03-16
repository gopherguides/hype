package hype

import (
	"encoding/json"
	"fmt"
)

type LinkCheckError struct {
	URL        string
	StatusCode int
	Err        error
}

func (e LinkCheckError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("link check failed for %q: status %d", e.URL, e.StatusCode)
	}
	if e.Err != nil {
		return fmt.Sprintf("link check failed for %q: %s", e.URL, e.Err)
	}
	return fmt.Sprintf("link check failed for %q", e.URL)
}

func (e LinkCheckError) Unwrap() error {
	return e.Err
}

func (e LinkCheckError) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":        "link_check_error",
		"url":         e.URL,
		"status_code": e.StatusCode,
	}
	if e.Err != nil {
		m["error"] = e.Err.Error()
	}
	return json.MarshalIndent(m, "", "  ")
}
