package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html/atom"
)

func Test_IsAtom(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		tag  Tag
		exp  bool
		want atom.Atom
	}{
		{name: "nil"},
		{
			name: "valid",
			want: atom.P,
			exp:  true,
			tag: &Element{
				Node: NewNode(htmx.ElementNode("p")),
			},
		},
		{
			name: "wrong atom",
			want: atom.Div,
			tag: &Element{
				Node: NewNode(htmx.ElementNode("p")),
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tt.exp, IsAtom(tt.tag, tt.want))
		})
	}
}

func Test_Tags_AllAtoms(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("big.html")
	r.NoError(err)

	act := doc.Children.AllAtom(atom.P)
	r.Len(act, 41)

	act = doc.Children.AllAtom(atom.Figure)
	r.Len(act, 23)

	act = doc.Children.AllAtom(atom.Textarea)
	r.Len(act, 0)
}

func Test_Tags_AllData(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("big.html")
	r.NoError(err)

	act := doc.Children.AllData("p")
	r.Len(act, 41)

	act = doc.Children.AllData("title")
	r.Len(act, 1)

	act = doc.Children.AllData("figure")
	r.Len(act, 23)

	act = doc.Children.AllData("textarea")
	r.Len(act, 0)
}

func Test_Tags_AllType(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("big.html")
	r.NoError(err)

	act := doc.Children.AllType(&Meta{})
	r.Len(act, 19)

	act = doc.Children.AllType(&InlineCode{})
	r.Len(act, 0)
}
