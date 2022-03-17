package hype

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseSnippets_Go(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	const path = "src/snippets.go"

	b, err := fs.ReadFile(testdata, path)
	r.NoError(err)

	p := testParser(t, testdata)

	snips, err := p.Snippets(path, b)
	r.NoError(err)
	r.NotEmpty(snips)
	r.Len(snips, 4)

	snip := snips["entertainer-funcs"]
	r.Equal("entertainer-funcs", snip.Name)
	r.Equal(path, snip.File)
	r.Equal("go", snip.Language)
	r.Equal(7, snip.Start)
	r.Equal(10, snip.End)

	exp := "\n\tName() string\n\tPerform(v Venue) error"
	act := snip.Content
	r.Equal(exp, act)
}

func Test_ParseSnippets_Ruby(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	const path = "src/snippets.rb"

	b, err := fs.ReadFile(testdata, path)
	r.NoError(err)

	p := testParser(t, testdata)

	snips, err := p.Snippets(path, b)
	r.NoError(err)
	r.NotEmpty(snips)
	r.Len(snips, 2)

	snip := snips["goodbye"]
	r.Equal("goodbye", snip.Name)
	r.Equal(path, snip.File)
	r.Equal("rb", snip.Language)
	r.Equal(7, snip.Start)
	r.Equal(11, snip.End)

	exp := "\ndef goodbye\n  puts \"Goodbye, World!\"\nend"
	act := snip.Content
	r.Equal(exp, act)
}
