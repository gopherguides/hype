package hype

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Include(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mod := `# Page 1

<include src="second/second.md"></include>`

	second := `# Second Page

<img src="assets/second.png"></img>`

	cab := fstest.MapFS{
		"module.md": &fstest.MapFile{
			Data: []byte(mod),
		},
		"second/second.md": &fstest.MapFile{
			Data: []byte(second),
		},
		"second/assets/second.png": &fstest.MapFile{},
	}

	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)
	r.NotNil(doc)

	exp := `<html><head></head><body><page>
<h1>Page 1</h1>
</page>
<page>
<h1>Second Page</h1>

<img src="second/assets/second.png"></img>
</page>

</body></html>`
	act := doc.String()

	// fmt.Println(act)
	r.Equal(exp, act)

}

func Test_Include_Nested(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mod := `# Page 1

<include src="second/second.md"></include>`

	second := `# Second Page

<img src="assets/second.png"></img>

<include src="third/third.md"></include>`

	third := `# Third Page

<img src="assets/third.png"></img>`

	cab := fstest.MapFS{
		"module.md": &fstest.MapFile{
			Data: []byte(mod),
		},
		"second/second.md": &fstest.MapFile{
			Data: []byte(second),
		},
		"second/assets/second.png": &fstest.MapFile{},
		"second/third/third.md": &fstest.MapFile{
			Data: []byte(third),
		},
		"second/third/assets/third.png": &fstest.MapFile{},
	}

	p := NewParser(cab)

	doc, err := p.ParseFile("module.md")
	r.NoError(err)
	r.NotNil(doc)

	exp := `<html><head></head><body><page>
<h1>Page 1</h1>
</page>
<page>
<h1>Second Page</h1>

<img src="second/assets/second.png"></img>
</page>
<page>
<h1>Third Page</h1>

<img src="second/third/assets/third.png"></img>
</page>


</body></html>`
	act := doc.String()

	// fmt.Println(act)
	r.Equal(exp, act)

}
