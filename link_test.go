package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Link_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	link := &Link{
		Element: NewEl("a", nil),
	}
	link.Nodes = append(link.Nodes, Text("This is a link"))

	err := link.Set("href", "https://example.com")
	r.NoError(err)

	testJSON(t, "link", link)
}
