package htmx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewAttributes(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	exp := Attributes{
		"id":  "1",
		"src": "foo.png",
	}

	node := AttrNode("img", exp)
	act := NewAttributes(node)
	r.Equal(exp, act)
}

func Test_Attributes_String(t *testing.T) {
	t.Parallel()

	ats := Attributes{
		"id":     "1",
		"quotey": `"This is a quote"`,
		"src":    "foo.png",
	}

	table := []struct {
		attrs Attributes
		exp   string
		name  string
	}{
		{name: "not empty", attrs: ats, exp: `id="1" quotey="\"This is a quote\"" src="foo.png"`},
		{name: "empty", exp: ""},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			act := tt.attrs.String()
			r.Equal(tt.exp, act)
		})
	}

}

func Test_Atributes_HasKeys(t *testing.T) {
	t.Parallel()

	ats := Attributes{
		"id":     "1",
		"quotey": `"This is a quote"`,
		"src":    "foo.png",
	}

	table := []struct {
		keys []string
		err  bool
		name string
	}{
		{name: "hit", keys: []string{"id", "quotey", "src"}},
		{name: "miss", keys: []string{"ID"}, err: true},
		{name: "empty keys", err: false},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			b := ats.HasKeys(tt.keys...)
			r.Equal(!tt.err, b)
		})
	}

}

func Test_Atributes_Matches(t *testing.T) {
	t.Parallel()

	ats := Attributes{
		"id":     "1",
		"quotey": `"This is a quote"`,
		"src":    "foo.png",
	}

	type query map[string]string

	table := []struct {
		query query
		err   bool
		name  string
	}{
		{name: "hit", query: query{"id": "1", "src": "foo.png"}},
		{name: "miss", query: query{"id": "2", "src": "foo.png"}, err: true},
		{name: "empty"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			b := ats.Matches(tt.query)
			r.Equal(!tt.err, b)
		})
	}

}
