package hype

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

	node := AttrNode(t, "img", exp)
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
