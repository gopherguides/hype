package hype

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PreParseError(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	ppe := PreParseError{
		Err: io.EOF,
	}

	r.True(ppe.As(&PreParseError{}), ppe)
	r.True(ppe.Is(ppe), ppe)
	r.True(ppe.Unwrap() == io.EOF, ppe)

	var pe PreParseError
	r.True(ppe.As(&pe), ppe)

	pe = PreParseError{}
	r.True(ppe.Is(pe), ppe)

	err := errors.Unwrap(ppe)
	r.Equal(io.EOF, err)
}
