package hype

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func byXYZCab(t testing.TB) *fstest.MapFS {
	t.Helper()

	mod := `# Page One

some text with a [link](http://example.com) inside of it.

<img src="foo.jpg" alt="foo" />

Here is some more text.

<img src="bar.jpg" alt="bar" />

and finally:

<code src="code.go"></code>
`

	cab := &fstest.MapFS{
		"hype.md": &fstest.MapFile{
			Data: []byte(mod),
		},
	}

	return cab
}

func byXYZDoc(t testing.TB) *Document {
	t.Helper()

	p := NewParser(byXYZCab(t))

	doc, err := p.ParseFile("hype.md")
	if err != nil {
		t.Fatal(err)
	}

	return doc
}

func Test_ByAtom(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := byXYZDoc(t)

	nodes := ByAtom(doc.Children(), "img")

	r.Len(nodes, 2)
}

func Test_ByAttrs(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := byXYZDoc(t)

	nodes := ByAttrs(doc.Children(), map[string]string{
		"src": "*",
	})

	r.Len(nodes, 3)
}

func Test_ByType(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := byXYZDoc(t)

	nodes := ByType[*SourceCode](doc.Children())

	r.Len(nodes, 1)
}
