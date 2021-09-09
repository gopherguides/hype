package hype

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseSnippets(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	b, err := fs.ReadFile(testdata, "src/snippets.go")
	r.NoError(err)

	snips, err := ParseSnippets("foo.go", b, nil)
	r.NoError(err)
	r.NotEmpty(snips)
	r.Len(snips, 4)

	snip := snips["entertainer-funcs"]
	r.Equal("entertainer-funcs", snip.Name)
	r.Equal("foo.go", snip.File)
	r.Equal("go", snip.Language)
	r.Equal(7, snip.Start)
	r.Equal(10, snip.End)

	exp := "\n\tName() string\n\tPerform(v Venue) error"
	act := snip.Content
	r.Equal(exp, act)
}
