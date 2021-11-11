package golang

import (
	"io"
	"os"
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

	if i.out != nil {
		w = io.MultiWriter(w, i.out)
	}

	i.out = w
	return i
}

func WithErr(i *StdIO, w io.Writer) *StdIO {
	if i == nil {
		i = &StdIO{}
	}

	if i.err != nil {
		w = io.MultiWriter(w, i.err)
	}

	i.err = w
	return i
}
