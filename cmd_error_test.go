package hype

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/markbates/clam"
	"github.com/stretchr/testify/require"
)

func Test_CmdError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	oce := CmdError{
		RunError: clam.RunError{
			Err: ExecuteError{
				Err: io.EOF,
			},
		},
	}

	wrapped := fmt.Errorf("error: %w", oce)

	r.True(oce.As(&CmdError{}), oce)
	r.True(oce.Is(oce), oce)
	r.True(oce.Unwrap() == io.EOF, oce)

	r.True(errors.As(wrapped, &CmdError{}), wrapped)
	r.True(errors.As(wrapped, &ExecuteError{}), wrapped)

	r.True(errors.Is(wrapped, CmdError{}), wrapped)
	r.True(errors.Is(wrapped, ExecuteError{}), wrapped)
	r.True(errors.Is(wrapped, io.EOF), wrapped)

	err := errors.Unwrap(oce)
	r.Equal(io.EOF, err)
}

func Test_CmdError_MarshalJSON(t *testing.T) {
	t.Parallel()

	ce := CmdError{
		RunError: clam.RunError{
			Args:   []string{"echo", "hello"},
			Env:    []string{"FOO=bar", "BAR=baz"},
			Err:    io.EOF,
			Exit:   1,
			Output: []byte("foo\nbar\nbaz\n"),
			Dir:    "/tmp",
		},
		Filename: "foo.go",
	}

	testJSON(t, "cmd_error", ce)
}
