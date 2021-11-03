package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_IsAtom(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		tag  Tag
		exp  bool
		want Atom
	}{
		{name: "nil"},
		{
			name: "valid",
			want: "p",
			exp:  true,
			tag: &Element{
				Node: NewNode(htmx.ElementNode("p")),
			},
		},
		{
			name: "wrong atom",
			want: "div",
			tag: &Element{
				Node: NewNode(htmx.ElementNode("p")),
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tt.exp, IsAtom(tt.tag, Atom(tt.want)))
		})
	}
}

func Test_Tags_AllAdam(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("big.html")
	r.NoError(err)

	act := doc.Children.ByAtom("p")
	r.Len(act, 41)

	act = doc.Children.ByAtom("title")
	r.Len(act, 1)

	act = doc.Children.ByAtom("figure")
	r.Len(act, 23)

	act = doc.Children.ByAtom("textarea")
	r.Len(act, 0)
}

func Test_Tags_AllType(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("big.html")
	r.NoError(err)

	act := doc.Children.ByType(&Meta{})
	r.Len(act, 19)

	act = doc.Children.ByType(&InlineCode{})
	r.Len(act, 0)
}

func Test_Tags_ByAttrs(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, testdata)

	doc, err := p.ParseFile("big.html")

	// rel="icon" type="image/png"

	r.NoError(err)
	r.NotNil(doc)

	table := []struct {
		name  string
		query Attributes
		exp   int
	}{
		{name: "hit", query: Attributes{
			"rel":  "icon",
			"type": "image/png",
		}, exp: 2},
		{name: "miss", query: Attributes{
			"rel":  "icon",
			"type": "image/jpeg",
		}, exp: 0},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			act := doc.Children.ByAttrs(tt.query)
			r.Len(act, tt.exp)
		})
	}

}
