package hype

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PreExecuteError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	pee := PreExecuteError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("error: %w", pee)

	r.True(pee.As(&PreExecuteError{}), pee)
	r.True(pee.Is(pee), pee)
	r.True(pee.Unwrap() == io.EOF, pee)

	var pe PreExecuteError
	r.True(errors.As(wrapped, &pe), wrapped)

	pe = PreExecuteError{}
	r.True(errors.Is(wrapped, pe), wrapped)

	err := errors.Unwrap(pee)
	r.Equal(io.EOF, err)
}
