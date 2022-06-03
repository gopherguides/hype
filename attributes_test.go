package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Language(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	good := &Attributes{}
	r.NoError(good.Set("language", "go"))

	class := &Attributes{}
	r.NoError(class.Set("class", "language-go"))

	short := &Attributes{}
	r.NoError(short.Set("lang", "go"))

	table := []struct {
		name string
		ats  *Attributes
		exp  string
	}{
		{name: "empty", ats: &Attributes{}, exp: "text"},
		{name: "good", ats: good, exp: "go"},
		{name: "prefix", ats: class, exp: "go"},
		{name: "short", ats: short, exp: "go"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			act := Language(tt.ats, "text")
			r.Equal(tt.exp, act)
		})
	}

}

func Test_AttrMatches(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	ats := &Attributes{}
	r.NoError(ats.Set("src", "foo.png"))

	table := []struct {
		name  string
		query map[string]string
		exp   bool
	}{
		{name: "empty map", query: map[string]string{}, exp: true},
		{name: "good", query: map[string]string{"src": "foo.png"}, exp: true},
		{name: "wildcard", query: map[string]string{"src": ".*"}, exp: true},
		{name: "bad", query: map[string]string{"src": "bar.png"}, exp: false},
		{name: "empty value", query: map[string]string{"src": ""}, exp: false},
		{name: "missing", query: map[string]string{"bar": "bar.png"}, exp: false},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			act := AttrMatches(ats, tc.query)
			r.Equal(tc.exp, act)

		})
	}

}
