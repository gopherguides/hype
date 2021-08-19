package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_json(t *testing.T) {
	t.Parallel()
	t.Skip()
	r := require.New(t)

	args := []string{"json", "html5.html"}

	rt, err := runtime(args)
	r.NoError(err)
	rt.Cab = os.DirFS("testdata")

	bb := &bytes.Buffer{}
	rt.Stdout = bb

	err = run(rt)
	r.NoError(err)

	act := strings.TrimSpace(bb.String())
	r.NotEmpty(act)

	r.Contains(act, `"fs":{"html5.html":{`)
	r.Contains(act, `"document":{"children":[`)
}

func Test_json_indent(t *testing.T) {
	t.Parallel()
	t.Skip()
	r := require.New(t)

	args := []string{"-i", "json", "html5.html"}

	rt, err := runtime(args)
	r.NoError(err)
	rt.Cab = os.DirFS("testdata")

	bb := &bytes.Buffer{}
	rt.Stdout = bb
	rt.SetOutput(bb)

	err = run(rt)

	act := strings.TrimSpace(bb.String())
	r.NotEmpty(act)

	r.NoError(err, act)

	r.Contains(act, "\"fs\": {\n")
	r.Contains(act, "{\n  \"document\"")
}

func Test_stream(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var args []string

	rt, err := runtime(args)
	r.NoError(err)
	r.NotNil(rt)

	in := `<p>Hi</hi>`
	rt.Stdin = strings.NewReader(in)

	bb := &bytes.Buffer{}
	rt.Stdout = bb

	err = stream(rt)
	r.NoError(err)

	act := strings.TrimSpace(bb.String())
	exp := "<html><head></head><body>\n<p></p><p>Hi</p>\n\n</body>\n</html>"

	r.Equal(exp, act)
}
