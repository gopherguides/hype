package hype

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Pages(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mod := `# Page 1

<include src="second/second.md"></include>

more text

<include src="third/third.md"></include>

adfadf`

	cab := fstest.MapFS{
		"module.md": &fstest.MapFile{
			Data: []byte(mod),
		},
		"second/second.md": &fstest.MapFile{
			Data: []byte(`# Second`),
		},
		"third/third.md": &fstest.MapFile{
			Data: []byte(`# Third`),
		},
	}

	p := NewParser(cab)

	doc, err := p.Parse(strings.NewReader(mod))
	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)
	exp := `<html><head></head><body><page>
<h1>Page 1</h1>
</page>
<page>
<h1>Second</h1>
</page>

<page>
<p>more text</p>
</page>
<page>
<h1>Third</h1>
</page>

<page>
<p>adfadf</p>
</page>
</body></html>`

	r.Equal(exp, act)

	pages := ByType[*Page](doc.Nodes)
	r.Len(pages, 5)

}

func Test_Pages_MarshalJSON(t *testing.T) {
	t.Parallel()

	p := &Page{
		Title:   "Page 1",
		Element: NewEl("page", nil),
	}
	p.Nodes = append(p.Nodes, Text("more text"))

	testJSON(t, "page", p)
}
