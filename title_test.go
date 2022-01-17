package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Document_Title(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		exp  string
	}{
		{name: "pages.md", exp: "First H1"},
		{name: "html5.html", exp: "A Basic HTML5 Template"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			doc := ParseFile(t, testdata, tt.name)
			r.NotNil(doc)

			r.Equal(tt.exp, doc.Title())
		})
	}

}

func Test_Document_SetTitle(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		exp  string
		err  bool
	}{
		{name: "pages.md", exp: "First H1"},
		{name: "html5.html", exp: "A Basic HTML5 Template"},
		{name: "notitle.md", exp: "Untitled", err: true},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			doc := ParseFile(t, testdata, tt.name)
			r.NotNil(doc)

			r.Equal(tt.exp, doc.Title())

			err := doc.SetTitle("New Title")
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal("New Title", doc.Title())
		})
	}

}
