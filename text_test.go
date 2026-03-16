package hype

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestText_MarshalJSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	tn := Text("hello")
	data, err := tn.MarshalJSON()
	r.NoError(err)

	var m map[string]any
	r.NoError(json.Unmarshal(data, &m))
	r.Equal("hello", m["text"])
	r.Contains(m["type"], "Text")
}

func TestText_Children(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	tn := Text("hello")
	r.Empty(tn.Children())
}

func TestText_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	r.Equal("hello", Text("hello").String())
	r.Equal("", Text("").String())
}

func TestText_MD(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	r.Equal("hello", Text("hello").MD())
}

func TestText_IsEmptyNode(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	r.True(Text("").IsEmptyNode())
	r.True(Text("  \t\n").IsEmptyNode())
	r.False(Text("hello").IsEmptyNode())
}
