package hype

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/markbates/clam"
	"github.com/markbates/hepa"
)

type CmdError struct {
	clam.RunError
	Filename string
}

func (ce CmdError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"args":     ce.Args,
		"env":      ce.Env,
		"error":    errForJSON(ce.Err),
		"exit":     ce.Exit,
		"filename": ce.Filename,
		"output":   string(ce.Output),
		"root":     ce.Dir,
		"type":     fmt.Sprintf("%T", ce),
	}

	p := hepa.Deep()

	env := make([]string, 0, len(ce.Env))
	for _, e := range ce.Env {
		b, _ := p.Clean(strings.NewReader(e))
		env = append(env, string(b))
	}
	mm["env"] = env

	return json.MarshalIndent(mm, "", "  ")
}

func (ce CmdError) Error() string {
	return toError(ce)
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
