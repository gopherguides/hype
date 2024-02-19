package hype

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PostExecuteError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	pee := PostExecuteError{
		Err: io.EOF,
	}

	wrapped := fmt.Errorf("error: %w", pee)

	r.True(pee.As(&PostExecuteError{}), pee)
	r.True(pee.Is(pee), pee)
	r.True(pee.Unwrap() == io.EOF, pee)

	var pe PostExecuteError
	r.True(errors.As(wrapped, &pe), wrapped)

	pe = PostExecuteError{}
	r.True(errors.Is(wrapped, pe), wrapped)

	err := errors.Unwrap(pee)
	r.Equal(io.EOF, err)

}
