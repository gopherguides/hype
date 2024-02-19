package hype

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	oce := ParseError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("error: %w", oce)

	r.True(oce.As(&ParseError{}), oce)
	r.True(oce.Is(oce), oce)
	r.True(oce.Unwrap() == io.EOF, oce)

	var ce ParseError
	r.True(errors.As(wrapped, &ce), wrapped)

	ce = ParseError{}
	r.True(errors.Is(wrapped, ce), wrapped)

	err := errors.Unwrap(oce)
	r.Equal(io.EOF, err)
}

func Test_ParseError_MarshalJSON(t *testing.T) {
	t.Parallel()

	pe := ParseError{
		Err:      io.EOF,
		Filename: "test.md",
		Root:     "root",
	}

	testJSON(t, "parse_error", pe)
}
