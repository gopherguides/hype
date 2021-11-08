package golang

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IO(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var i *StdIO

	r.Equal(os.Stderr, i.Err())
	r.Equal(os.Stdin, i.In())
	r.Equal(os.Stdout, i.Out())

	bb := &bytes.Buffer{}

	i = WithIn(i, bb)
	r.Equal(bb, i.In())

	i = WithOut(i, bb)
	r.Equal(bb, i.Out())

	i = WithErr(i, bb)
	r.Equal(bb, i.Err())
}
