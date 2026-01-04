package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_YouTube_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	yt := &YouTube{
		Element: NewEl("youtube", nil),
	}

	err := yt.Set("id", "dQw4w9WgXcQ")
	r.NoError(err)

	err = yt.Set("title", "Rick Astley - Never Gonna Give You Up")
	r.NoError(err)

	testJSON(t, "youtube", yt)
}

func Test_YouTube_String(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	yt := &YouTube{
		Element: NewEl("youtube", nil),
	}

	err := yt.Set("id", "dQw4w9WgXcQ")
	r.NoError(err)

	err = yt.Set("title", "Test Video")
	r.NoError(err)

	got := yt.String()

	r.Contains(got, `<div class="youtube-embed">`)
	r.Contains(got, `src="https://www.youtube.com/embed/dQw4w9WgXcQ"`)
	r.Contains(got, `title="Test Video"`)
	r.Contains(got, `allowfullscreen`)
	r.Contains(got, `</iframe>`)
	r.Contains(got, `</div>`)
}

func Test_YouTube_String_DefaultTitle(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	yt := &YouTube{
		Element: NewEl("youtube", nil),
	}

	err := yt.Set("id", "dQw4w9WgXcQ")
	r.NoError(err)

	got := yt.String()

	r.Contains(got, `title="YouTube video player"`)
}

func Test_YouTube_ValidationError_MissingID(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	el := NewEl("youtube", nil)
	_, err := NewYouTube(el)

	r.Error(err)
	r.Contains(err.Error(), "id")
}

func Test_YouTube_ValidationError_InvalidID(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		id   string
	}{
		{name: "too short", id: "abc"},
		{name: "too long", id: "abcdefghijklmnop"},
		{name: "invalid chars", id: "abc!@#$%^&*"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			el := NewEl("youtube", nil)
			err := el.Set("id", tc.id)
			r.NoError(err)

			_, err = NewYouTube(el)
			r.Error(err)
			r.Contains(err.Error(), "invalid YouTube video ID")
		})
	}
}

func Test_YouTube_VideoID(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	yt := &YouTube{
		Element: NewEl("youtube", nil),
	}

	err := yt.Set("id", "dQw4w9WgXcQ")
	r.NoError(err)

	id, err := yt.VideoID()
	r.NoError(err)
	r.Equal("dQw4w9WgXcQ", id)
}

func Test_YouTube_Title(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	yt := &YouTube{
		Element: NewEl("youtube", nil),
	}

	err := yt.Set("id", "dQw4w9WgXcQ")
	r.NoError(err)

	err = yt.Set("title", "My Video")
	r.NoError(err)

	r.Equal("My Video", yt.Title())
}

func Test_YouTube_MD(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	yt := &YouTube{
		Element: NewEl("youtube", nil),
	}

	err := yt.Set("id", "dQw4w9WgXcQ")
	r.NoError(err)

	got := yt.MD()
	r.Contains(got, `<div class="youtube-embed">`)
	r.Contains(got, `src="https://www.youtube.com/embed/dQw4w9WgXcQ"`)
}
