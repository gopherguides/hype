package binding

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrPath_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	e := ErrPath("/some/path")
	r.Contains(e.Error(), "/some/path")
	r.Contains(e.Error(), "could not parse section from")
}

func TestErrPath_Is(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	e := ErrPath("/path1")
	r.True(e.Is(ErrPath("/path2")))
	r.False(e.Is(errors.New("other error")))
}
