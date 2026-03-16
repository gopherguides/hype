package blog

import (
	"bytes"
	"html/template"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSEO_OpenGraphTags(t *testing.T) {
	t.Parallel()

	t.Run("full struct", func(t *testing.T) {
		r := require.New(t)
		s := SEO{
			Title:       "My Title",
			Description: "A description",
			URL:         "https://example.com",
			Image:       "https://example.com/img.png",
			Type:        "article",
		}
		tags := s.OpenGraphTags()
		r.Contains(tags, `og:title`)
		r.Contains(tags, `My Title`)
		r.Contains(tags, `og:description`)
		r.Contains(tags, `og:url`)
		r.Contains(tags, `og:image`)
		r.Contains(tags, `og:type`)
		r.Contains(tags, `article`)
	})

	t.Run("minimal struct defaults type to website", func(t *testing.T) {
		r := require.New(t)
		s := SEO{Title: "Title Only"}
		tags := s.OpenGraphTags()
		r.Contains(tags, `og:title`)
		r.Contains(tags, `website`)
		r.NotContains(tags, `og:description`)
		r.NotContains(tags, `og:url`)
		r.NotContains(tags, `og:image`)
	})

	t.Run("escapes special characters", func(t *testing.T) {
		r := require.New(t)
		s := SEO{Title: `Title & "quotes"`}
		tags := s.OpenGraphTags()
		r.Contains(tags, `Title &amp; &quot;quotes&quot;`)
	})
}

func TestSEO_TwitterCardTags(t *testing.T) {
	t.Parallel()

	t.Run("full struct", func(t *testing.T) {
		r := require.New(t)
		s := SEO{
			Title:       "My Title",
			Description: "A description",
			Image:       "https://example.com/img.png",
			TwitterCard: "summary",
			TwitterSite: "@example",
		}
		tags := s.TwitterCardTags()
		r.Contains(tags, `twitter:card`)
		r.Contains(tags, `summary`)
		r.Contains(tags, `twitter:title`)
		r.Contains(tags, `twitter:description`)
		r.Contains(tags, `twitter:image`)
		r.Contains(tags, `twitter:site`)
		r.Contains(tags, `@example`)
	})

	t.Run("defaults card to summary_large_image", func(t *testing.T) {
		r := require.New(t)
		s := SEO{Title: "Title"}
		tags := s.TwitterCardTags()
		r.Contains(tags, `summary_large_image`)
	})

	t.Run("omits empty fields", func(t *testing.T) {
		r := require.New(t)
		s := SEO{Title: "Title"}
		tags := s.TwitterCardTags()
		r.NotContains(tags, `twitter:description`)
		r.NotContains(tags, `twitter:image`)
		r.NotContains(tags, `twitter:site`)
	})
}

func TestSEO_JSONLD(t *testing.T) {
	t.Parallel()

	t.Run("returns empty for non-article", func(t *testing.T) {
		r := require.New(t)
		s := SEO{Title: "Title", Type: "website"}
		r.Empty(s.JSONLD())
	})

	t.Run("returns empty when type is empty", func(t *testing.T) {
		r := require.New(t)
		s := SEO{Title: "Title"}
		r.Empty(s.JSONLD())
	})

	t.Run("returns script tag for article", func(t *testing.T) {
		r := require.New(t)
		s := SEO{
			Title:     "My Article",
			Type:      "article",
			Author:    "John Doe",
			Published: "2024-01-15",
		}
		jsonld := s.JSONLD()
		r.Contains(jsonld, `application/ld+json`)
		r.Contains(jsonld, `"My Article"`)
		r.Contains(jsonld, `"John Doe"`)
		r.Contains(jsonld, `"2024-01-15"`)
		r.Contains(jsonld, `schema.org`)
	})
}

func TestEscapeAttr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"ampersand", "a&b", "a&amp;b"},
		{"quote", `a"b`, "a&quot;b"},
		{"less than", "a<b", "a&lt;b"},
		{"greater than", "a>b", "a&gt;b"},
		{"empty", "", ""},
		{"no special chars", "hello world", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tt.out, escapeAttr(tt.in))
		})
	}
}

func renderSEOPartial(t *testing.T, data PageData) string {
	t.Helper()

	seoBytes, err := os.ReadFile("templates/partials/seo.html")
	require.NoError(t, err)

	tmplText := string(seoBytes) + `{{template "partials/seo" .}}`
	tmpl, err := template.New("test").Parse(tmplText)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	return buf.String()
}

func TestSEOArticleTags(t *testing.T) {
	published := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	article := Article{
		Title:          "Getting Started with Hype",
		Slug:           "getting-started",
		Published:      published,
		Author:         "Gopher Guides",
		AuthorTwitter:  "@caborundrum",
		Tags:           []string{"tutorial", "getting-started", "hype"},
		SEODescription: "Learn how to install Hype",
		OGImage:        "/images/getting-started-og.png",
	}

	cfg := Config{
		Title:       "Hype",
		Description: "A blog powered by hype",
		BaseURL:     "https://hypemd.dev",
		Author: Author{
			Name:    "Gopher Guides",
			Twitter: "@hype_markdown",
		},
		SEO: SEOConfig{
			OGImage:     "/images/og-default.png",
			TwitterCard: "summary_large_image",
			SiteName:    "Hype",
		},
	}

	data := PageData{
		Config:  cfg,
		Article: &article,
	}

	out := renderSEOPartial(t, data)

	t.Run("og:site_name from config", func(t *testing.T) {
		require.Contains(t, out, `og:site_name" content="Hype"`)
	})

	t.Run("article:published_time", func(t *testing.T) {
		require.Contains(t, out, `article:published_time" content="2026-03-15T00:00:00Z"`)
	})

	t.Run("article:author", func(t *testing.T) {
		require.Contains(t, out, `article:author" content="Gopher Guides"`)
	})

	t.Run("article:tag for each tag", func(t *testing.T) {
		require.Contains(t, out, `article:tag" content="tutorial"`)
		require.Contains(t, out, `article:tag" content="getting-started"`)
		require.Contains(t, out, `article:tag" content="hype"`)
	})

	t.Run("twitter:site from config", func(t *testing.T) {
		require.Contains(t, out, `twitter:site" content="@hype_markdown"`)
	})

	t.Run("twitter:creator from article author_twitter", func(t *testing.T) {
		require.Contains(t, out, `twitter:creator" content="@caborundrum"`)
	})

	t.Run("per-article og:image", func(t *testing.T) {
		require.Contains(t, out, `og:image" content="https://hypemd.dev/images/getting-started-og.png"`)
	})

	t.Run("per-article twitter:image", func(t *testing.T) {
		require.Contains(t, out, `twitter:image" content="https://hypemd.dev/images/getting-started-og.png"`)
	})
}

func TestSEOArticleFallbacks(t *testing.T) {
	published := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	article := Article{
		Title:     "Minimal Article",
		Slug:      "minimal",
		Published: published,
	}

	cfg := Config{
		Title:       "My Blog",
		Description: "Default description",
		BaseURL:     "https://example.com",
		Author: Author{
			Name:    "Default Author",
			Twitter: "@default_handle",
		},
		SEO: SEOConfig{
			OGImage:     "/images/default.png",
			TwitterCard: "summary_large_image",
		},
	}

	data := PageData{
		Config:  cfg,
		Article: &article,
	}

	out := renderSEOPartial(t, data)

	t.Run("og:site_name falls back to Config.Title", func(t *testing.T) {
		require.Contains(t, out, `og:site_name" content="My Blog"`)
	})

	t.Run("article:author falls back to config author", func(t *testing.T) {
		require.Contains(t, out, `article:author" content="Default Author"`)
	})

	t.Run("twitter:creator falls back to config twitter", func(t *testing.T) {
		require.Contains(t, out, `twitter:creator" content="@default_handle"`)
	})

	t.Run("og:image falls back to config default", func(t *testing.T) {
		require.Contains(t, out, `og:image" content="https://example.com/images/default.png"`)
	})

	t.Run("no article:tag when tags empty", func(t *testing.T) {
		require.NotContains(t, out, `article:tag`)
	})
}

func TestSEOListPage(t *testing.T) {
	cfg := Config{
		Title:       "Hype",
		Description: "A blog powered by hype",
		BaseURL:     "https://hypemd.dev",
		Author: Author{
			Name:    "Gopher Guides",
			Twitter: "@hype_markdown",
		},
		SEO: SEOConfig{
			OGImage:     "/images/og-default.png",
			TwitterCard: "summary_large_image",
			SiteName:    "Hype Blog",
		},
	}

	data := PageData{
		Config:  cfg,
		Article: nil,
	}

	out := renderSEOPartial(t, data)

	t.Run("og:type is website", func(t *testing.T) {
		require.Contains(t, out, `og:type" content="website"`)
	})

	t.Run("og:site_name from SEO config", func(t *testing.T) {
		require.Contains(t, out, `og:site_name" content="Hype Blog"`)
	})

	t.Run("twitter:site from config", func(t *testing.T) {
		require.Contains(t, out, `twitter:site" content="@hype_markdown"`)
	})

	t.Run("no twitter:creator on list page", func(t *testing.T) {
		require.NotContains(t, out, `twitter:creator`)
	})

	t.Run("no article-specific tags on list page", func(t *testing.T) {
		require.NotContains(t, out, `article:published_time`)
		require.NotContains(t, out, `article:author`)
		require.NotContains(t, out, `article:tag`)
	})
}

func TestSEOBackwardCompatibility(t *testing.T) {
	cfg := Config{
		Title:       "My Blog",
		Description: "A blog",
		BaseURL:     "https://example.com",
		SEO: SEOConfig{
			TwitterCard: "summary_large_image",
		},
	}

	t.Run("article with no new fields", func(t *testing.T) {
		article := Article{
			Title:     "Test",
			Slug:      "test",
			Published: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		data := PageData{Config: cfg, Article: &article}
		out := renderSEOPartial(t, data)

		require.Contains(t, out, `og:type" content="article"`)
		require.Contains(t, out, `og:title" content="Test"`)
		require.NotContains(t, out, `twitter:site`)
		require.NotContains(t, out, `twitter:creator`)
		require.Contains(t, out, `og:site_name" content="My Blog"`)
	})

	t.Run("list page with no new fields", func(t *testing.T) {
		data := PageData{Config: cfg, Article: nil}
		out := renderSEOPartial(t, data)

		require.Contains(t, out, `og:type" content="website"`)
		require.NotContains(t, out, `twitter:site`)
		_ = strings.Contains(out, `og:site_name`)
	})
}
