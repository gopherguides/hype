package hype

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ExtractMeta_Basic(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/json-meta/basic"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	doc, err := p.ParseFile("hype.md")
	r.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = doc.Execute(ctx)
	r.NoError(err)

	meta := ExtractMeta(doc)

	r.Equal("Getting Started with Go", meta.Title)
	r.Equal("getting-started-with-go", meta.Slug)

	r.Equal("Cory LaNou", meta.Metadata["author"])
	r.Equal("2026-02-08", meta.Metadata["published"])
	r.Equal("go, tutorial, beginner", meta.Metadata["tags"])

	r.Len(meta.Headings, 4)
	r.Equal(1, meta.Headings[0].Level)
	r.Equal("Getting Started with Go", meta.Headings[0].Text)
	r.Equal(2, meta.Headings[1].Level)
	r.Equal("Installation", meta.Headings[1].Text)
	r.Equal(2, meta.Headings[2].Level)
	r.Equal("First Program", meta.Headings[2].Text)
	r.Equal(3, meta.Headings[3].Level)
	r.Equal("Hello World", meta.Headings[3].Text)

	r.Len(meta.CodeSnippets, 1)
	r.Equal("go", meta.CodeSnippets[0].Language)
	r.Nil(meta.CodeSnippets[0].Src)

	r.Empty(meta.Includes)
	r.Empty(meta.Images)

	r.Greater(meta.WordCount, 0)
	r.Greater(meta.ReadingTimeMinutes, 0)

	b, err := json.Marshal(meta)
	r.NoError(err)
	r.True(json.Valid(b))
}

func Test_ExtractMeta_MissingFields(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/json-meta/minimal"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	doc, err := p.ParseFile("hype.md")
	r.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = doc.Execute(ctx)
	r.NoError(err)

	meta := ExtractMeta(doc)

	r.Equal("Minimal Document", meta.Title)
	r.Equal("minimal-document", meta.Slug)
	r.Empty(meta.Metadata)
	r.Len(meta.Headings, 1)
	r.Empty(meta.CodeSnippets)
	r.Empty(meta.Includes)
	r.Empty(meta.Images)
	r.Greater(meta.WordCount, 0)
}

func Test_ExtractTOC_Hierarchy(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/json-meta/complex-toc"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	doc, err := p.ParseFile("hype.md")
	r.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = doc.Execute(ctx)
	r.NoError(err)

	toc := ExtractTOC(doc)

	r.Equal("Top Level", toc.Title)
	r.Len(toc.TOC, 1)

	top := toc.TOC[0]
	r.Equal(1, top.Level)
	r.Equal("Top Level", top.Text)
	r.Len(top.Children, 3)

	ch1 := top.Children[0]
	r.Equal(2, ch1.Level)
	r.Equal("Chapter One", ch1.Text)
	r.Len(ch1.Children, 2)

	s1a := ch1.Children[0]
	r.Equal(3, s1a.Level)
	r.Equal("Section One A", s1a.Text)
	r.Empty(s1a.Children)

	s1b := ch1.Children[1]
	r.Equal(3, s1b.Level)
	r.Equal("Section One B", s1b.Text)
	r.Len(s1b.Children, 1)

	sub := s1b.Children[0]
	r.Equal(4, sub.Level)
	r.Equal("Subsection One B i", sub.Text)

	ch2 := top.Children[1]
	r.Equal(2, ch2.Level)
	r.Equal("Chapter Two", ch2.Text)
	r.Len(ch2.Children, 1)

	ch3 := top.Children[2]
	r.Equal(2, ch3.Level)
	r.Equal("Chapter Three", ch3.Text)
	r.Empty(ch3.Children)

	b, err := json.Marshal(toc)
	r.NoError(err)
	r.True(json.Valid(b))
}

func Test_Slugify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"Getting Started with Go", "getting-started-with-go"},
		{"Hello, World!", "hello-world"},
		{"  Spaces  Everywhere  ", "spaces-everywhere"},
		{"already-slugged", "already-slugged"},
		{"UPPERCASE", "uppercase"},
		{"multiple---dashes", "multiple-dashes"},
		{"special@#chars$%", "special-chars"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tc.want, slugify(tc.input))
		})
	}
}
