package hype

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Source_MimeType(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		src  string
		exp  string
	}{
		{name: "known ext mime", src: "example.go", exp: "text/x-go; charset=utf-8"},
		{name: "missing ext", src: "example", exp: "text/plain"},
		{name: "unknown ext mime", src: "example.unknown", exp: "text/plain"},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			src := Source(tt.src)

			act := src.MimeType()
			r.Equal(tt.exp, act)
		})
	}

}

func Test_Source_Lang(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		src  string
		exp  string
	}{
		{name: "file ext", src: "example.md", exp: "md"},
		{name: "missing ext", src: "example", exp: "plain"},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			src := Source(tt.src)

			act := src.Lang()
			r.Equal(tt.exp, act)
		})
	}

}

func Test_Source_Scheme(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		exp  string
		src  string
		err  bool
	}{
		{name: "http source", exp: "http", src: "http://example.com"},
		{name: "https source", exp: "https", src: "https://example.com"},
		{name: "file source", exp: "file", src: "file:///tmp/example.md"},
		{name: "empty scheme", exp: "file", src: "/tmp/example.md"},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			src := Source(tt.src)

			act, err := src.Scheme()
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal(tt.exp, act)
		})
	}

}

func Test_Source_StatFile(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	ref, err := fs.Stat(testdata, "pages.md")

	r.NoError(err)

	table := []struct {
		name string
		src  string
		cab  fs.FS
		err  bool
		exp  fs.FileInfo
	}{
		{name: "file source schemeless", src: "pages.md", cab: testdata, exp: ref},
		{name: "file source with scheme", src: "file://pages.md", cab: testdata, exp: ref},
		{name: "unknown file", src: "unknown", cab: testdata, err: true},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			src := Source(tt.src)

			act, err := src.StatFile(tt.cab)
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.Equal(tt.exp, act)

		})
	}
}

type tripper func(*http.Request) (*http.Response, error)

func (t tripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return t(r)
}

func Test_Source_StatHTTP(t *testing.T) {
	t.Parallel()

	const fn = "test.md"

	r := require.New(t)

	now, err := time.Parse(time.Kitchen, time.Kitchen)
	r.NoError(err)

	body := []byte("Hello, World")

	testdata := fstest.MapFS{
		fn: {
			Data:    body,
			Mode:    fs.ModeIrregular,
			ModTime: now,
		},
	}

	ref, err := fs.Stat(testdata, fn)
	r.NoError(err)

	client := http.DefaultClient

	client.Transport = tripper(func(r *http.Request) (*http.Response, error) {
		path := r.URL.Path

		if path != "/"+fn {
			return &http.Response{
				StatusCode: 404,
			}, nil
		}

		res := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header: http.Header{
				"Content-Type": []string{"text/markdown"},
				"Content-Length": []string{
					fmt.Sprintf("%d", ref.Size()),
				},
				"Last-Modified": []string{
					ref.ModTime().Format(http.TimeFormat),
				},
			},
		}

		return res, nil
	})

	table := []struct {
		name string
		src  string
		err  bool
		exp  fs.FileInfo
	}{
		{name: "good http", src: "http://example.com/" + fn, exp: ref},
		{name: "good https", src: "https://example.com/" + fn, exp: ref},
		{name: "404", src: "https://example.com/unknown", exp: ref, err: true},
		{name: "file", src: "file://" + fn, err: true},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			src := Source(tt.src)

			act, err := src.StatHTTP(client)
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)

			exp := tt.exp
			r.Equal(exp.Size(), act.Size())
			r.Equal(exp.ModTime(), act.ModTime())
			r.Equal(exp.Mode(), act.Mode())

		})
	}
}
