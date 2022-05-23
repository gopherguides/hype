package hype

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Snippets(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	sm := &Snippets{}

	sm.Add(".go", "// %s")

	_, ok := sm.Get("main.go")
	r.False(ok)
}

func Test_Parse_Snippets(t *testing.T) {
	t.Parallel()

	const goexp = "func Goodbye() {\n\tfmt.Println(\"Goodbye, World!\")\n}"
	const rbexp = "def goodbye\n  puts \"Goodbye, World!\"\nend"
	const jsexp = "function goodbye() {\n    console.log('Goodbye, World!');\n}"
	const htmlexp = "<p>Goodbye World</p>"

	table := []struct {
		path  string
		lang  string
		start int
		end   int
		exp   string
	}{
		{path: "snippets.go", lang: "go", start: 12, end: 17, exp: goexp},
		{path: "snippets.rb", lang: "rb", start: 7, end: 11, exp: rbexp},
		{path: "snippets.js", lang: "js", start: 7, end: 11, exp: jsexp},
		{path: "snippets.html", lang: "html", start: 16, end: 18, exp: htmlexp},
	}

	for _, tc := range table {
		t.Run(tc.path, func(t *testing.T) {
			r := require.New(t)

			b, err := fs.ReadFile(os.DirFS("testdata/snippets"), tc.path)
			r.NoError(err)

			sm := &Snippets{}

			snips, err := sm.Parse(tc.path, b)

			r.NoError(err)
			r.NotEmpty(snips)
			r.Len(snips, 2)

			snip := snips["goodbye"]
			r.Equal("goodbye", snip.Name)
			r.Equal(tc.path, snip.File)
			r.Equal(tc.lang, snip.Lang)
			r.Equal(tc.start, snip.Start)
			r.Equal(tc.end, snip.End)

			act := snip.Content
			r.Equal(tc.exp, act)
		})
	}

}

func Test_ParseSnippets_Unclosed(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	src := `package demo
	// snippet: example
	func Hello() {
	}
	`

	sm := &Snippets{}

	_, err := sm.Parse("demo.go", []byte(src))
	r.Error(err)

}

func Test_ParseSnippets_Duplicate(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	src := `package demo
	// snippet: example
	func Hello() {
	}
	// snippet: example
	// snippet: example
	func Goodbye() {
	}
	// snippet: example
	`

	sm := &Snippets{}

	_, err := sm.Parse("demo.go", []byte(src))
	r.Error(err)
}
