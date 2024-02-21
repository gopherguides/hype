package hype

import (
	"context"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Document_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		r := require.New(t)
		p := testParser(t, "testdata/doc/execution/success")

		doc, err := p.ParseFile("module.md")
		r.NoError(err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err = doc.Execute(ctx)
		r.NoError(err)

		act := doc.String()
		act = strings.TrimSpace(act)

		exp := "<html><head></head><body><page>\n<h1>Command</h1>\n\n<cmd exec=\"echo 'Hello World'\"><pre><code class=\"language-shell\" language=\"shell\">$ echo Hello World\n\nHello World</code></pre></cmd>\n</page>\n</body></html>"

		r.Equal(exp, act)
	})

	t.Run("failure", func(t *testing.T) {
		r := require.New(t)
		p := testParser(t, "testdata/doc/execution/failure")

		doc, err := p.ParseFile("module.md")
		r.NoError(err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err = doc.Execute(ctx)
		r.Error(err)

		_, ok := err.(ExecuteError)
		r.True(ok, err)
		r.True(errors.Is(err, ExecuteError{}), err)
	})

}

func Test_Document_MD(t *testing.T) {
	t.Skip("TODO: fix this test")
	t.Parallel()
	r := require.New(t)

	root := "testdata/doc/to_md"

	p := testParser(t, root)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	r.NoError(err)

	s := doc.MD()
	act := string(s)
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	b, err := fs.ReadFile(p.FS, "module.gold")
	r.NoError(err)

	exp := string(b)
	exp = strings.TrimSpace(exp)

	r.Equal(exp, act)
}

func Test_Document_MarshalJSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/doc/snippets")
	p.DocIDGen = func() (string, error) {
		return "1", nil
	}

	ctx := context.Background()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	r.NoError(err)

	testJSON(t, "document", doc)
}

func Test_Document_Pages(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/doc/pages")

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	pages, err := doc.Pages()
	r.NoError(err)

	r.Len(pages, 3)
}

func Test_Document_Pages_NoPages(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/doc/simple")

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	pages, err := doc.Pages()
	r.NoError(err)

	r.Len(pages, 1)
}
