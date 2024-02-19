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
			Err: io.EOF,
		},
	}

	wrapped := fmt.Errorf("error: %w", oce)

	r.True(oce.As(&CmdError{}), oce)
	r.True(oce.Is(oce), oce)
	r.True(oce.Unwrap() == io.EOF, oce)

	var ce CmdError
	r.True(errors.As(wrapped, &ce), wrapped)

	ce = CmdError{}
	r.True(errors.Is(wrapped, ce), wrapped)

	err := errors.Unwrap(oce)
	r.Equal(io.EOF, err)
}
