package blog

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestArticle_IsPublished(t *testing.T) {
	t.Parallel()

	t.Run("past date returns true", func(t *testing.T) {
		r := require.New(t)
		a := Article{Published: time.Now().Add(-24 * time.Hour)}
		r.True(a.IsPublished())
	})

	t.Run("future date returns false", func(t *testing.T) {
		r := require.New(t)
		a := Article{Published: time.Now().Add(24 * time.Hour)}
		r.False(a.IsPublished())
	})
}

func TestArticle_URL(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	a := Article{Slug: "my-post"}
	r.Equal("/my-post/", a.URL())
}

func TestArticle_FormattedDate(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	a := Article{Published: time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC)}
	r.Equal("March 15, 2024", a.FormattedDate())
}

func TestArticle_ISODate(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	a := Article{Published: time.Date(2024, time.March, 15, 10, 30, 0, 0, time.UTC)}
	iso := a.ISODate()
	r.Contains(iso, "2024-03-15")
	r.Contains(iso, "T10:30:00")
}

func TestCalculateReadingTime(t *testing.T) {
	t.Parallel()

	t.Run("200 words is 1 minute", func(t *testing.T) {
		r := require.New(t)
		words := strings.Repeat("word ", 200)
		r.Equal(1, calculateReadingTime(words))
	})

	t.Run("400 words is 2 minutes", func(t *testing.T) {
		r := require.New(t)
		words := strings.Repeat("word ", 400)
		r.Equal(2, calculateReadingTime(words))
	})

	t.Run("1 word rounds up to 1 minute", func(t *testing.T) {
		r := require.New(t)
		r.Equal(1, calculateReadingTime("hello"))
	})

	t.Run("empty is 0 minutes", func(t *testing.T) {
		r := require.New(t)
		r.Equal(0, calculateReadingTime(""))
	})

	t.Run("strips HTML before counting", func(t *testing.T) {
		r := require.New(t)
		html := "<p>" + strings.Repeat("word ", 200) + "</p>"
		r.Equal(1, calculateReadingTime(html))
	})
}

func TestStripHTML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"removes tags", "<p>hello</p>", "hello"},
		{"nested tags", "<div><p>hello</p></div>", "hello"},
		{"preserves text", "no tags here", "no tags here"},
		{"empty", "", ""},
		{"self closing", "a<br/>b", "ab"},
		{"attributes", `<a href="x">link</a>`, "link"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tt.out, stripHTML(tt.in))
		})
	}
}

func TestStripDetailsBlocks(t *testing.T) {
	t.Parallel()

	t.Run("removes details block", func(t *testing.T) {
		r := require.New(t)
		input := "before<details><summary>Click</summary>Hidden</details>after"
		r.Equal("beforeafter", stripDetailsBlocks(input))
	})

	t.Run("preserves non-details content", func(t *testing.T) {
		r := require.New(t)
		r.Equal("hello world", stripDetailsBlocks("hello world"))
	})

	t.Run("removes multiple details blocks", func(t *testing.T) {
		r := require.New(t)
		input := "a<details>1</details>b<details>2</details>c"
		r.Equal("abc", stripDetailsBlocks(input))
	})
}
