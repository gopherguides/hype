package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/markbates/clam"
	"github.com/markbates/hepa"
)

type CmdError struct {
	clam.RunError `json:"clam_error,omitempty"`
	Filename      string `json:"filename,omitempty"`
}

func (ce CmdError) MarshalJSON() ([]byte, error) {
	mm := map[string]any{
		"args":     ce.Args,
		"env":      ce.Env,
		"error":    ce.Err,
		"exit":     ce.Exit,
		"filename": ce.Filename,
		"output":   string(ce.Output),
		"root":     ce.Dir,
		"type":     fmt.Sprintf("%T", ce),
	}

	if _, ok := ce.Err.(json.Marshaler); !ok && ce.Err != nil {
		mm["err"] = ce.Err.Error()
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
	b, _ := json.MarshalIndent(ce, "", "  ")
	return string(b)
}

func (ce CmdError) Unwrap() error {
	type Unwrapper interface {
		Unwrap() error
	}

	if _, ok := ce.Err.(Unwrapper); ok {
		return errors.Unwrap(ce.Err)
	}

	return ce.Err
}

func (ce CmdError) Is(target error) bool {
	if ce.Err == nil {
		return false
	}

	return errors.Is(ce.Err, target)
}

func (ce CmdError) As(target any) bool {
	return errors.As(ce.Err, target)
}
