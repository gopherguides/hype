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

//go:embed templates/* templates/partials/*
var defaultTemplates embed.FS

type Renderer struct {
	blog         *Blog
	singleTmpl   *template.Template
	listTmpl     *template.Template
	xmlTemplates *texttemplate.Template
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

	baseTmpl := template.New("").Funcs(htmlFuncMap)
	baseTmpl, _ = baseTmpl.ParseFS(defaultTemplates, "templates/partials/*.html")
	baseTmpl, _ = baseTmpl.ParseFS(defaultTemplates, "templates/baseof.html")

	singleTmpl, _ := baseTmpl.Clone()
	singleTmpl, _ = singleTmpl.ParseFS(defaultTemplates, "templates/single.html")

	listTmpl, _ := baseTmpl.Clone()
	listTmpl, _ = listTmpl.ParseFS(defaultTemplates, "templates/list.html")

	xmlTmpl := texttemplate.New("").Funcs(xmlFuncMap)
	xmlTmpl, _ = xmlTmpl.ParseFS(defaultTemplates, "templates/*.xml")

	return &Renderer{
		blog:         b,
		singleTmpl:   singleTmpl,
		listTmpl:     listTmpl,
		xmlTemplates: xmlTmpl,
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
	if err := r.listTmpl.ExecuteTemplate(&buf, "baseof", data); err != nil {
		return fmt.Errorf("failed to execute list template: %w", err)
	}

	indexPath := filepath.Join(outDir, "index.html")
	return os.WriteFile(indexPath, buf.Bytes(), 0644)
}

func (r *Renderer) RenderArticle(outDir string, article Article) error {
	data := r.newPageData()
	data.Article = &article

	var buf bytes.Buffer
	if err := r.singleTmpl.ExecuteTemplate(&buf, "baseof", data); err != nil {
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
