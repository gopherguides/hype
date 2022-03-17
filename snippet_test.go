package hype

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse_Snippets(t *testing.T) {
	t.Parallel()

	const goexp = "\nfunc Goodbye() {\n\tfmt.Println(\"Goodbye, World!\")\n}\n"
	const rbexp = "\ndef goodbye\n  puts \"Goodbye, World!\"\nend"
	const jsexp = "\nfunction goodbye() {\n    console.log('Goodbye, World!');\n}"

	table := []struct {
		path  string
		lang  string
		start int
		end   int
		exp   string
	}{
		{path: "src/snippets.go", lang: "go", start: 12, end: 17, exp: goexp},
		{path: "src/snippets.rb", lang: "rb", start: 7, end: 11, exp: rbexp},
		{path: "src/snippets.js", lang: "js", start: 7, end: 11, exp: jsexp},
	}

	p := testParser(t, testdata)

	for _, tc := range table {
		t.Run(tc.path, func(t *testing.T) {
			r := require.New(t)

			b, err := fs.ReadFile(testdata, tc.path)
			r.NoError(err)

			snips, err := p.Snippets(tc.path, b)
			r.NoError(err)
			r.NotEmpty(snips)
			r.Len(snips, 2)

			snip := snips["goodbye"]
			r.Equal("goodbye", snip.Name)
			r.Equal(tc.path, snip.File)
			r.Equal(tc.lang, snip.Language)
			r.Equal(tc.start, snip.Start)
			r.Equal(tc.end, snip.End)

			act := snip.Content
			r.Equal(tc.exp, act)
		})
	}

}

func Test_ParseSnippets_Go(t *testing.T) {
	t.Parallel()
	t.Skip()
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
	t.Skip()
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
