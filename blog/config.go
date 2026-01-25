package blog

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	BaseURL     string    `yaml:"baseURL"`
	Author      Author    `yaml:"author"`
	Theme       string    `yaml:"theme"`
	Highlight   Highlight `yaml:"highlight"`
	SEO         SEOConfig `yaml:"seo"`
	ContentDir  string    `yaml:"contentDir"`
	OutputDir   string    `yaml:"outputDir"`
}

type Author struct {
	Name    string `yaml:"name"`
	Email   string `yaml:"email"`
	Twitter string `yaml:"twitter"`
}

type Highlight struct {
	Style       string `yaml:"style"`
	LineNumbers bool   `yaml:"lineNumbers"`
}

type SEOConfig struct {
	OGImage     string `yaml:"ogImage"`
	TwitterCard string `yaml:"twitterCard"`
}

func DefaultConfig() Config {
	return Config{
		Title:       "My Blog",
		Description: "A blog powered by hype",
		Theme:       "suspended",
		Highlight: Highlight{
			Style:       "monokai",
			LineNumbers: false,
		},
		SEO: SEOConfig{
			TwitterCard: "summary_large_image",
		},
		ContentDir: "content",
		OutputDir:  "public",
	}
}

func LoadConfig(fsys fs.FS) (Config, error) {
	cfg := DefaultConfig()

	data, err := fs.ReadFile(fsys, "config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			data, err = fs.ReadFile(fsys, "config.yml")
		}
		if err != nil {
			return cfg, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.ContentDir == "" {
		cfg.ContentDir = "content"
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = "public"
	}

	return cfg, nil
}

func LoadConfigFromPath(path string) (Config, error) {
	dir := filepath.Dir(path)
	fsys := os.DirFS(dir)
	return LoadConfig(fsys)
}
