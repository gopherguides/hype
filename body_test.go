package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Body_AsPage(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	body := &Body{
		Element: NewEl("body", nil),
	}

	p := body.AsPage()
	r.NotNil(p)
	r.Equal(body.Element, p.Element)

}
