package hype

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Var(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	data := map[string]any{
		"id": 1,
	}

	fn, err := NewVarParserFn(data)
	r.NoError(err)
	r.NotNil(fn)

	tcs := []struct {
		name string
		key  string
		exp  string
		err  bool
	}{
		{name: "valid", key: "id", exp: "1"},
		{name: "unknown key", key: "404", err: true},
		{name: "empty key", key: "", err: true},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			el := NewEl("var", nil)
			r.NotNil(el)
			el.Nodes = Nodes{Text(tc.key)}

			nodes, err := fn(nil, el)
			if tc.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(nodes)
			r.NotEmpty(nodes)

			doc := &Document{
				Nodes: nodes,
			}

			err = doc.Execute(context.Background())
			r.NoError(err)

			r.Equal(tc.exp, doc.String())
		})
	}

}

func Test_Var_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var v *Var
	r.Equal("<var></var>", v.String())

	v = &Var{}
	r.Equal("<var></var>", v.String())

	v.Element = &Element{}
	r.Equal("<var></var>", v.String())

	v.Nodes = append(v.Nodes, Text("id"))
	r.Equal("<var>id</var>", v.String())

	v.value = 1
	r.Equal("1", v.String())
}
