package hype

import (
	"testing"
	"testing/fstest"

	"github.com/gopherguides/hype/htmltest"
	"github.com/stretchr/testify/require"
)

func Test_Parser_NewMeta(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	p := testParser(t, fstest.MapFS{})

	_, err := p.NewMeta(nil)
	r.Error(err)

	const key = `props`
	const val = `yo!`

	table := []struct {
		ats  Attributes
		err  bool
		name string
	}{
		{name: "both keys", ats: Attributes{"property": key, "name": key, "content": val}, err: true},
		{name: "missing content", ats: Attributes{"property": key}, err: true},
		{name: "missing key", err: true},
		{name: "name key", ats: Attributes{"name": key, "content": val}},
		{name: "nil", err: true},
		{name: "property key", ats: Attributes{"property": key, "content": val}},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			n := htmltest.AttrNode(t, "meta", tt.ats)

			node := NewNode(n)

			m, err := p.NewMeta(node)
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(m)
			r.Equal(key, m.Key)
			r.Equal(val, m.Val)
		})
	}

}

func Test_Metas_Value(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	ms := Metas{
		{Key: "src", Val: "foo.png"},
	}

	v, ok := ms.Value("src")
	r.True(ok)
	r.Equal("foo.png", v)

	v, ok = ms.Value("404")
	r.False(ok)
	r.Empty(v)
}
