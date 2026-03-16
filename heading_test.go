package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Slug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"Introduction", "introduction"},
		{"My Document", "my-document"},
		{"Hello World!", "hello-world"},
		{"  spaces  ", "spaces"},
		{"Special @#$ Characters", "special-characters"},
		{"already-slug", "already-slug"},
		{"UPPER CASE", "upper-case"},
		{"multiple   spaces", "multiple-spaces"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tt.want, Slug(tt.in))
		})
	}
}

func Test_UniqueSlug(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	seen := map[string]int{}

	r.Equal("example", UniqueSlug("Example", seen))
	r.Equal("example-1", UniqueSlug("Example", seen))
	r.Equal("example-2", UniqueSlug("Example", seen))
	r.Equal("other", UniqueSlug("Other", seen))
	r.Equal("other-1", UniqueSlug("Other", seen))
}

func Test_UniqueSlug_EmptyFallback(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	seen := map[string]int{}

	r.Equal("heading", UniqueSlug("!!!", seen))
	r.Equal("heading-1", UniqueSlug("###", seen))
	r.Equal("heading-2", UniqueSlug("", seen))
}

func Test_UniqueSlug_PreSeededCollision(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	seen := map[string]int{}
	seen["example"] = 1
	seen["example-1"] = 1

	r.Equal("example-2", UniqueSlug("Example", seen))
}

func Test_NewHeading(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	h, err := NewHeading(&Element{
		HTMLNode: &html.Node{
			Type: html.ElementNode,
			Data: "h1",
		},
	})

	r.NoError(err)
	r.Equal("<h1></h1>", h.String())
	r.Equal(1, h.Level())
}

func Test_Heading_MarshalJSON(t *testing.T) {
	t.Parallel()

	h := &Heading{
		Element: NewEl("h1", nil),
		level:   1,
	}
	h.Nodes = append(h.Nodes, Text("This is a heading"))

	testJSON(t, "heading", h)
}
