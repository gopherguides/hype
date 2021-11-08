package golang

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type StdIO struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func (i *StdIO) In() io.Reader {
	if i == nil || i.in == nil {
		return os.Stdin
	}
	return i.in
}

func (i *StdIO) Out() io.Writer {
	if i == nil || i.out == nil {
		return os.Stdout
	}
	return i.out
}

func (i *StdIO) Err() io.Writer {
	if i == nil || i.err == nil {
		return os.Stderr
	}
	return i.err
}

func WithIn(i *StdIO, r io.Reader) *StdIO {
	if i == nil {
		i = &StdIO{}
	}
	i.in = r
	return i
}

func WithOut(i *StdIO, w io.Writer) *StdIO {
	if i == nil {
		i = &StdIO{}
	}
	i.out = w
	return i
}

func WithErr(i *StdIO, w io.Writer) *StdIO {
	if i == nil {
		i = &StdIO{}
	}
	i.err = w
	return i
}

func execute(ctx context.Context, std *StdIO, args ...string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	c := exec.CommandContext(ctx, "go", args...)

	bb := &bytes.Buffer{}

	c.Stdout = bb
	c.Stderr = std.Err()
	c.Stdin = std.In()

	err := c.Run()

	if err != nil {
		cargs := strings.Join(c.Args, " ")
		return fmt.Errorf("$ %s: %w", cargs, err)
	}

	s := html.EscapeString(bb.String())
	io.Copy(std.Out(), strings.NewReader(s))

	return nil
}
