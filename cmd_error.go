package hype

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/markbates/clam"
)

type CmdError struct {
	clam.RunError
	Filename string
}

func (ce CmdError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"args":     ce.Args,
		"error":    errForJSON(ce.Err),
		"exit":     ce.Exit,
		"filename": ce.Filename,
		"output":   string(ce.Output),
		"root":     ce.Dir,
		"type":     toType(ce),
	}

	return json.MarshalIndent(mm, "", "  ")
}

func (ce CmdError) Error() string {
	sb := &strings.Builder{}

	var lines []string

	fp := filepath.Join(ce.Dir, ce.Filename)
	if len(fp) > 0 {
		lines = append(lines, fmt.Sprintf("filepath: %s", fp))
	}

	if len(ce.Args) > 0 {
		lines = append(lines, fmt.Sprintf("cmd: $ %s", strings.Join(ce.Args, " ")))
	}

	if ce.Exit != 0 {
		lines = append(lines, fmt.Sprintf("exit: %d", ce.Exit))
	}

	if ce.Err != nil {
		lines = append(lines, fmt.Sprintf("error: %s", ce.Err))
	}

	sb.WriteString(strings.Join(lines, "\n"))

	s := sb.String()

	return strings.TrimSpace(s)
}

func (ce CmdError) As(target any) bool {
	ex, ok := target.(*CmdError)
	if !ok {
		return ce.RunError.As(target)
	}

	(*ex) = ce
	return true
}

func (ce CmdError) Is(target error) bool {
	if _, ok := target.(CmdError); ok {
		return true
	}

	return ce.RunError.Is(target)
}
