package hype

import (
	"testing"

	"github.com/gopherguides/hype/atomx"
	"github.com/stretchr/testify/require"
)

func Test_Parse_Fragment(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	in := `<a href="foo.zip" target="_blank">text</a>`
	frag := []byte(in)

	tags, err := p.ParseFragment(frag)
	r.NoError(err)
	r.Len(tags, 1)

	tag := tags[0]
	r.Equal(atomx.A, tag.Atom())
	r.Equal(in, tag.String())

}
