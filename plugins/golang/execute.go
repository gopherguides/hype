package golang

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"os/exec"
	"strings"
)

type Runner struct {
	IO   *StdIO
	Args []string
}

func NewRunner(args ...string) *Runner {
	return &Runner{
		Args: args,
	}
}

func (r *Runner) CMD(ctx context.Context) *exec.Cmd {
	c := exec.CommandContext(ctx, "go", r.Args...)
	c.Stdout = r.IO.Out()
	c.Stderr = r.IO.Err()
	c.Stdin = r.IO.In()
	return c
}

func (r *Runner) String() string {
	return fmt.Sprintf("go %s", strings.Join(r.Args, " "))
}

func (r *Runner) Run(ctx context.Context) error {
	c := r.CMD(ctx)
	return c.Run()
}

func execute(ctx context.Context, std *StdIO, args ...string) error {
	r := NewRunner(args...)
	r.IO = std

	c := r.CMD(ctx)

	bb := &bytes.Buffer{}

	c.Stdout = bb

	err := c.Run()

	if err != nil {
		return fmt.Errorf("$ %s: %w", r, err)
	}

	s := html.EscapeString(bb.String())
	io.Copy(std.Out(), strings.NewReader(s))

	return nil
}
