package hype

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Markdown_UnknownAtom(t *testing.T) {
	t.Skip()
	t.Parallel()
	r := require.New(t)

	root := "testdata/markdown/unknown-atom"

	p := testParser(t, root)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	exp := `<a for="io#EOF" href="https://pkg.go.dev/io#EOF" target="_blank"><code>io.EOF</code></a>`

	act := doc.String()
	act = strings.TrimSpace(act)

	// fmt.Println(act)
	r.Equal(act, exp)
}

func Test_Markdown(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	fn := Markdown()

	in := strings.NewReader(`# Hello World`)

	out, err := fn(&Parser{}, in)
	r.NoError(err)

	b, err := io.ReadAll(out)
	r.NoError(err)

	act := string(b)
	act = strings.TrimSpace(act)
	exp := "<page>\n<h1>Hello World</h1>\n</page>"

	r.Equal(exp, act)
}

func Test_Markdown_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	fn := Markdown()

	_, err := fn(nil, nil)
	r.Error(err)

	_, err = fn(nil, brokenReader{})
	r.Error(err)
}
