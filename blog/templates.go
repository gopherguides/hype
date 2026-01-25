package blog

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	texttemplate "text/template"
)

//go:embed templates/*
var defaultTemplates embed.FS

type Renderer struct {
	blog         *Blog
	htmlTemplates *template.Template
	xmlTemplates  *texttemplate.Template
}

func NewRenderer(b *Blog) *Renderer {
	htmlFuncMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
		"join": strings.Join,
	}

	xmlFuncMap := texttemplate.FuncMap{
		"join": strings.Join,
	}

	htmlTmpl := template.New("").Funcs(htmlFuncMap)
	htmlTmpl, _ = htmlTmpl.ParseFS(defaultTemplates, "templates/*.html")

	xmlTmpl := texttemplate.New("").Funcs(xmlFuncMap)
	xmlTmpl, _ = xmlTmpl.ParseFS(defaultTemplates, "templates/*.xml")

	return &Renderer{
		blog:         b,
		htmlTemplates: htmlTmpl,
		xmlTemplates:  xmlTmpl,
	}
}

type PageData struct {
	Config      Config
	Article     *Article
	Articles    []Article
	HighlightCSS template.CSS
	CurrentYear int
}

func (r *Renderer) newPageData() PageData {
	return PageData{
		Config:       r.blog.Config,
		Articles:     r.blog.Articles,
		HighlightCSS: template.CSS(r.blog.Highlighter.CSS()),
		CurrentYear:  2026,
	}
}

func (r *Renderer) RenderIndex(outDir string) error {
	data := r.newPageData()

	var buf bytes.Buffer
	if err := r.htmlTemplates.ExecuteTemplate(&buf, "list.html", data); err != nil {
		return fmt.Errorf("failed to execute list template: %w", err)
	}

	indexPath := filepath.Join(outDir, "index.html")
	return os.WriteFile(indexPath, buf.Bytes(), 0644)
}

func (r *Renderer) RenderArticle(outDir string, article Article) error {
	data := r.newPageData()
	data.Article = &article

	var buf bytes.Buffer
	if err := r.htmlTemplates.ExecuteTemplate(&buf, "single.html", data); err != nil {
		return fmt.Errorf("failed to execute single template: %w", err)
	}

	articleDir := filepath.Join(outDir, article.Slug)
	if err := os.MkdirAll(articleDir, 0755); err != nil {
		return err
	}

	articlePath := filepath.Join(articleDir, "index.html")
	return os.WriteFile(articlePath, buf.Bytes(), 0644)
}

func (r *Renderer) RenderRSS(outDir string) error {
	data := r.newPageData()

	var buf bytes.Buffer
	if err := r.xmlTemplates.ExecuteTemplate(&buf, "rss.xml", data); err != nil {
		return fmt.Errorf("failed to execute RSS template: %w", err)
	}

	rssPath := filepath.Join(outDir, "rss.xml")
	return os.WriteFile(rssPath, buf.Bytes(), 0644)
}

func (r *Renderer) RenderSitemap(outDir string) error {
	data := r.newPageData()

	var buf bytes.Buffer
	if err := r.xmlTemplates.ExecuteTemplate(&buf, "sitemap.xml", data); err != nil {
		return fmt.Errorf("failed to execute sitemap template: %w", err)
	}

	sitemapPath := filepath.Join(outDir, "sitemap.xml")
	return os.WriteFile(sitemapPath, buf.Bytes(), 0644)
}

func (r *Renderer) RenderRobots(outDir string) error {
	content := fmt.Sprintf("User-agent: *\nAllow: /\n\nSitemap: %s/sitemap.xml\n", r.blog.Config.BaseURL)
	robotsPath := filepath.Join(outDir, "robots.txt")
	return os.WriteFile(robotsPath, []byte(content), 0644)
}
