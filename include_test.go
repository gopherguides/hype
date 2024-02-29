package hype

import (
	"context"
	"errors"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Include_Parse_Errors(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		root     string
		filename string
	}{
		{
			root:     "testdata/includes/broken",
			filename: "hype.md",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.root, func(t *testing.T) {
			r := require.New(t)
			p := testParser(t, tc.root)

			ctx := context.Background()

			_, err := p.ParseExecuteFile(ctx, "hype.md")
			r.Error(err)

			var ee ParseError
			r.True(errors.As(err, &ee))

			r.Equal(tc.filename, ee.Filename)
			r.Equal(tc.root, ee.Root)

		})
	}

}

func Test_Include_Cmd_Errors(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		root     string
		filename string
	}{
		{
			root:     "testdata/includes/toplevel",
			filename: "hype.md",
		},
		{
			root:     "testdata/includes/sublevel",
			filename: "below/b.md",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.root, func(t *testing.T) {
			r := require.New(t)
			p := testParser(t, tc.root)

			ctx := context.Background()

			_, err := p.ParseExecuteFile(ctx, "hype.md")
			r.Error(err)

			var ee ExecuteError
			r.True(errors.As(err, &ee))

			r.Equal(tc.filename, ee.Filename)
			r.Equal(tc.root, ee.Root)

			var ce CmdError
			r.True(errors.As(ee.Err, &ce))

			r.Equal(tc.filename, ce.Filename)
			r.Equal(-1, ce.Exit)
		})
	}

}

func Test_Include(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mod := `# Page 1

<include src="second/second.md"></include>`

	second := `# Second Page

<img src="assets/second.png"></img>`

	cab := fstest.MapFS{
		"hype.md": &fstest.MapFile{
			Data: []byte(mod),
		},
		"second/second.md": &fstest.MapFile{
			Data: []byte(second),
		},
		"second/assets/second.png": &fstest.MapFile{},
	}

	p := NewParser(cab)

	doc, err := p.ParseFile("hype.md")
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
		"hype.md": &fstest.MapFile{
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

	doc, err := p.ParseFile("hype.md")
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

func Test_Include_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	inc := &Include{
		Element: NewEl("include", nil),
		dir:     "testdata/includes",
	}

	err := inc.Set("src", "sub.md")
	r.NoError(err)

	testJSON(t, "include", inc)
}
