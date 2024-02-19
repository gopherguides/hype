package hype

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExecuteError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	oce := ExecuteError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("error: %w", oce)

	r.True(oce.As(&ExecuteError{}), oce)
	r.True(oce.Is(oce), oce)
	r.True(oce.Unwrap() == io.EOF, oce)

	var ce ExecuteError
	r.True(errors.As(wrapped, &ce), wrapped)

	ce = ExecuteError{}
	r.True(errors.Is(wrapped, ce), wrapped)

	err := errors.Unwrap(oce)
	r.Equal(io.EOF, err)
}
