package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Env_Getenv(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	e := &Env{}
	r.Equal("", e.Getenv("unknown"))

	v := e.Getenv("PATH")
	r.NotEmpty(v)
	r.NotEqual("bar", v)

	e.Setenv("PATH", "bar")
	v = e.Getenv("PATH")
	r.NotEmpty(v)
	r.Equal("bar", v)
}
