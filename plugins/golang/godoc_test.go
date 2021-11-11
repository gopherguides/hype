package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_Godoc_CleanFlags(t *testing.T) {
	t.Parallel()

	empty := []string{}
	blank := []string{"-a", " ", "-b", "", "-c"}
	spaces := []string{" -a ", "\t-b", "-c\n"}
	good := []string{"-a", "-b", "-c"}

	table := []struct {
		name  string
		flags []string
		exp   []string
	}{
		{name: "empty", flags: empty, exp: empty},
		{name: "blank", flags: blank, exp: good},
		{name: "spaces", flags: spaces, exp: good},
		{name: "good", flags: good, exp: good},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			act := CleanFlags(tt.flags...)

			r.Equal(tt.exp, act)
		})
	}

}

func Test_Godoc(t *testing.T) {
	t.Parallel()
	t.Skip()
	r := require.New(t)

	root, err := CachePath()
	r.NoError(err)
	// os.RemoveAll(root)

	table := []struct {
		root  string
		src   string
		flags []string
		exp   string
	}{
		{src: "encoding/json#Decoder", flags: []string{"-short"}, exp: filepath.Join("encoding", "json", "Decoder.short.godoc")},
		{src: "github.com/gobuffalo/flect#Ident", exp: filepath.Join("github.com", "gobuffalo", "flect", "Ident.godoc")},
	}

	for _, tt := range table {
		t.Run(fmt.Sprintf("%s/%s", tt.src, tt.flags), func(t *testing.T) {
			r := require.New(t)

			g := &Godoc{}

			act := g.key(tt.src, tt.flags...)

			r.Equal(tt.exp, act)

			node := hype.NewNode(htmx.AttrNode("godoc", htmx.Attributes{
				"src":   tt.src,
				"flags": strings.Join(tt.flags, ","),
			}))

			gd, err := NewGodoc(node)
			r.NoError(err)
			r.NotNil(gd)

			_, err = os.Stat(filepath.Join(root, tt.exp))
			r.NoError(err)
		})
	}

}

func Test_Parser_Godoc(t *testing.T) {
	t.Parallel()
	t.Skip()
	r := require.New(t)

	var root = filepath.Join("testdata", "godoc.parser.test")

	cab := os.DirFS(root)

	p := testParser(t, cab, root)

	p.SetCustomTag("godoc", func(node *hype.Node) (hype.Tag, error) {
		return NewGodoc(node)
	})

	f, err := os.Open(filepath.Join("testdata", "godoc.md"))
	r.NoError(err)
	defer f.Close()

	doc, err := p.ParseReader(f)
	r.NoError(err)
	r.NotNil(doc)

	exp := expOutput
	act := doc.String()

	// fmt.Println(act)
	r.Equal(exp, act)

}

const expOutput = `<html><head></head><body>
# Godoc Tag

Some text

<godoc src="io#EOF"><pre><code language="godoc" class="language-godoc">$ go doc io.EOF

// Go Version:		1.17.3
// Documentation:	<a href="https://pkg.go.dev/io#EOF" target="_blank">https://pkg.go.dev/io#EOF</a>

package io // import "io"

var EOF = errors.New("EOF")
    EOF is the error returned by Read when no more input is available. (Read
    must return EOF itself, not an error wrapping EOF, because callers will test
    for EOF using ==.) Functions should return EOF only to signal a graceful end
    of input. If the EOF occurs unexpectedly in a structured data stream, the
    appropriate error is either ErrUnexpectedEOF or some other error giving more
    detail.

</code></pre></godoc>

some more text

<godoc flags="-short" src="golang.org/x/net/html#TokenType.String"><pre><code language="godoc" class="language-godoc">$ go doc -short golang.org/x/net/html.TokenType.String

// Go Version:		1.17.3
// Documentation:	<a href="https://pkg.go.dev/golang.org/x/net/html#TokenType.String" target="_blank">https://pkg.go.dev/golang.org/x/net/html#TokenType.String</a>

func (t TokenType) String() string
    String returns a string representation of the TokenType.

</code></pre></godoc>

even more text
</body>
</html>`
