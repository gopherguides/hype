package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Var(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	data := map[string]any{
		"id": 1,
	}

	p := testParser(t, "testdata")
	err := p.Vars.BulkSet(data)
	r.NoError(err)

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

			nodes, err := NewVarNodes(p, el)
			if tc.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(nodes)
			r.NotEmpty(nodes)

			doc := &Document{
				Nodes:  nodes,
				Parser: p,
			}

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

	v = &Var{
		Element: NewEl("var", nil),
	}

	r.Equal("<var></var>", v.String())

	v.Nodes = append(v.Nodes, Text("id"))
	r.Equal("<var>id</var>", v.String())

	v.Value = 1
	r.Equal("1", v.String())
}

func Test_Var_MarshalJSON(t *testing.T) {
	t.Parallel()

	v := &Var{
		Element: NewEl("var", nil),
		Key:     "id",
		Value:   1,
	}
	v.Nodes = append(v.Nodes, Text("id"))

	testJSON(t, "var", v)

}
