package blog

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Blog struct {
	Config      Config
	Articles    []Article
	Highlighter *Highlighter
	fsys        fs.FS
	root        string
}

func New(root string) (*Blog, error) {
	fsys := os.DirFS(root)

	cfg, err := LoadConfig(fsys)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	h := NewHighlighter(cfg.Highlight.Style, cfg.Highlight.LineNumbers)

	return &Blog{
		Config:      cfg,
		Highlighter: h,
		fsys:        fsys,
		root:        root,
	}, nil
}

func (b *Blog) Discover(ctx context.Context) error {
	contentDir := b.Config.ContentDir

	contentFS, err := fs.Sub(b.fsys, contentDir)
	if err != nil {
		return fmt.Errorf("failed to access content directory: %w", err)
	}

	entries, err := fs.ReadDir(contentFS, ".")
	if err != nil {
		return fmt.Errorf("failed to read content directory: %w", err)
	}

	parser := NewArticleParser(contentFS, b.Highlighter, b.Config.BaseURL)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}

		article, err := parser.ParseArticle(ctx, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", name, err)
			continue
		}

		if !article.IsPublished() {
			continue
		}

		b.Articles = append(b.Articles, article)
	}

	sort.Slice(b.Articles, func(i, j int) bool {
		return b.Articles[i].Published.After(b.Articles[j].Published)
	})

	return nil
}

func (b *Blog) Build(ctx context.Context) error {
	if err := b.Discover(ctx); err != nil {
		return err
	}

	outDir := filepath.Join(b.root, b.Config.OutputDir)

	if err := validateOutputDir(b.root, outDir); err != nil {
		return err
	}

	if err := os.RemoveAll(outDir); err != nil {
		return fmt.Errorf("failed to clean output directory: %w", err)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	renderer := NewRenderer(b)

	if err := renderer.RenderIndex(outDir); err != nil {
		return fmt.Errorf("failed to render index: %w", err)
	}

	for _, article := range b.Articles {
		if err := renderer.RenderArticle(outDir, article); err != nil {
			return fmt.Errorf("failed to render %s: %w", article.Slug, err)
		}
	}

	if err := renderer.RenderRSS(outDir); err != nil {
		return fmt.Errorf("failed to render RSS: %w", err)
	}

	if err := renderer.RenderSitemap(outDir); err != nil {
		return fmt.Errorf("failed to render sitemap: %w", err)
	}

	if err := renderer.RenderRobots(outDir); err != nil {
		return fmt.Errorf("failed to render robots.txt: %w", err)
	}

	staticDir := filepath.Join(b.root, "static")
	if _, err := os.Stat(staticDir); err == nil {
		if err := copyDir(staticDir, outDir); err != nil {
			return fmt.Errorf("failed to copy static files: %w", err)
		}
	}

	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}

func validateOutputDir(root, outDir string) error {
	if outDir == "" {
		return fmt.Errorf("output directory cannot be empty")
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("failed to resolve root path: %w", err)
	}

	absOut, err := filepath.Abs(outDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	if absOut == absRoot {
		return fmt.Errorf("output directory cannot be the project root")
	}

	if !strings.HasPrefix(absOut, absRoot+string(filepath.Separator)) {
		return fmt.Errorf("output directory must be inside the project root")
	}

	return nil
}
