package blog

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cfg := DefaultConfig()

	r.Equal("My Blog", cfg.Title)
	r.Equal("A blog powered by hype", cfg.Description)
	r.Equal("suspended", cfg.Theme)
	r.Equal("monokai", cfg.Highlight.Style)
	r.False(cfg.Highlight.LineNumbers)
	r.Equal("summary_large_image", cfg.SEO.TwitterCard)
	r.Equal("content", cfg.ContentDir)
	r.Equal("public", cfg.OutputDir)
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	t.Run("parses config.yaml", func(t *testing.T) {
		r := require.New(t)
		fsys := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte("title: Test Blog\ndescription: A test\nbaseURL: https://test.com\ncontentDir: posts\noutputDir: dist\n"),
			},
		}
		cfg, err := LoadConfig(fsys)
		r.NoError(err)
		r.Equal("Test Blog", cfg.Title)
		r.Equal("A test", cfg.Description)
		r.Equal("https://test.com", cfg.BaseURL)
		r.Equal("posts", cfg.ContentDir)
		r.Equal("dist", cfg.OutputDir)
	})

	t.Run("falls back to config.yml", func(t *testing.T) {
		r := require.New(t)
		fsys := fstest.MapFS{
			"config.yml": &fstest.MapFile{
				Data: []byte("title: YML Blog\n"),
			},
		}
		cfg, err := LoadConfig(fsys)
		r.NoError(err)
		r.Equal("YML Blog", cfg.Title)
	})

	t.Run("missing config returns error", func(t *testing.T) {
		r := require.New(t)
		fsys := fstest.MapFS{}
		_, err := LoadConfig(fsys)
		r.Error(err)
		r.Contains(err.Error(), "failed to read config file")
	})

	t.Run("malformed YAML returns error", func(t *testing.T) {
		r := require.New(t)
		fsys := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte("title: [invalid yaml\n"),
			},
		}
		_, err := LoadConfig(fsys)
		r.Error(err)
		r.Contains(err.Error(), "failed to parse config")
	})

	t.Run("empty contentDir defaults", func(t *testing.T) {
		r := require.New(t)
		fsys := fstest.MapFS{
			"config.yaml": &fstest.MapFile{
				Data: []byte("title: Blog\ncontentDir: \"\"\noutputDir: \"\"\n"),
			},
		}
		cfg, err := LoadConfig(fsys)
		r.NoError(err)
		r.Equal("content", cfg.ContentDir)
		r.Equal("public", cfg.OutputDir)
	})
}
