package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FindTitle(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	doc := byXYZDoc(t)

	title := FindTitle(doc.Children())

	r.Equal("Page One", title)
}
