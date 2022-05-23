package hype

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Markdown(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	fn := Markdown()

	in := strings.NewReader(`# Hello World`)

	out, err := fn(&Parser{}, in)
	r.NoError(err)

	b, err := io.ReadAll(out)
	r.NoError(err)

	act := string(b)
	act = strings.TrimSpace(act)
	exp := "<page>\n<h1>Hello World</h1>\n</page>"

	r.Equal(exp, act)
}

func Test_Markdown_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	fn := Markdown()

	_, err := fn(nil, nil)
	r.Error(err)

	_, err = fn(nil, brokenReader{})
	r.Error(err)
}
