package hype

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_GenerateToC(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/toc/basic")
	p.Section = 42

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	act := doc.String()
	act = strings.TrimSpace(act)

	r.Contains(act, `<nav class="toc">`)
	r.Contains(act, `href="#table-of-contents"`)
	r.Contains(act, `href="#aaa"`)
	r.Contains(act, `href="#bbb"`)
	r.Contains(act, `href="#bbb1"`)
	r.Contains(act, `<h1 id="table-of-contents">`)
	r.Contains(act, `<h2 id="aaa">`)
}

func Test_GenerateToC_Depth(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/toc/depth")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	act := doc.String()

	r.Contains(act, `href="#introduction"`)
	r.Contains(act, `href="#details"`)
	r.NotContains(act, `href="#sub-detail"`)
}

func Test_GenerateToC_NoRoot(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/toc/noroot")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	act := doc.String()

	r.NotContains(act, `href="#my-document"`)
	r.Contains(act, `href="#introduction"`)
	r.Contains(act, `href="#conclusion"`)
}

func Test_GenerateToC_Duplicates(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/toc/duplicates")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	act := doc.String()

	r.Contains(act, `href="#example"`)
	r.Contains(act, `href="#example-1"`)
	r.Contains(act, `href="#example-2"`)
	r.Contains(act, `id="example"`)
	r.Contains(act, `id="example-1"`)
	r.Contains(act, `id="example-2"`)
}

func Test_ToC_MD(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/toc/basic")
	p.Section = 42

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	md := doc.MD()

	r.Contains(md, "[Table of Contents](#table-of-contents)")
	r.Contains(md, "[AAA](#aaa)")
}

func Test_ToC_MarshalJSON(t *testing.T) {
	t.Parallel()

	toc := &ToC{
		Element: NewEl("toc", nil),
		Depth:   6,
		Root:    true,
	}

	testJSON(t, "toc", toc)
}
