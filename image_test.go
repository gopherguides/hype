package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Image_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	i := &Image{
		Element: NewEl("img", nil),
	}

	err := i.Set("src", "https://example.com/image.jpg")
	r.NoError(err)

	testJSON(t, "image", i)
}
