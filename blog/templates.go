package blog

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	texttemplate "text/template"
)

//go:embed themes/*/theme.yaml themes/*/layouts/* themes/*/layouts/**/*
var embeddedThemes embed.FS

var builtinThemes = []string{"suspended", "developer", "cards"}

type Renderer struct {
	blog         *Blog
	singleTmpl   *template.Template
	listTmpl     *template.Template
	xmlTemplates *texttemplate.Template
}

type templateLayer struct {
	name string
	fsys fs.FS
	base string
}

func NewRenderer(b *Blog) (*Renderer, error) {
	htmlFuncMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
		"join": strings.Join,
		"slice": func(s string, start, end int) string {
			if start < 0 {
				start = 0
			}
			if end > len(s) {
				end = len(s)
			}
			if start >= end || start >= len(s) {
				return ""
			}
			return s[start:end]
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}

	xmlFuncMap := texttemplate.FuncMap{
		"join": strings.Join,
	}

	layers := buildTemplateLayers(b)

	baseTmpl := template.New("").Funcs(htmlFuncMap)
	baseTmpl, err := parsePartialsFromLayers(baseTmpl, layers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse partials: %w", err)
	}
	baseTmpl, err = parseFileFromLayers(baseTmpl, layers, "_default/baseof.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse baseof: %w", err)
	}

	singleTmpl, err := baseTmpl.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone template: %w", err)
	}
	singleTmpl, err = parseFileFromLayers(singleTmpl, layers, "_default/single.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse single: %w", err)
	}

	listTmpl, err := baseTmpl.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone template: %w", err)
	}
	listTmpl, err = parseFileFromLayers(listTmpl, layers, "_default/list.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse list: %w", err)
	}

	xmlTmpl := texttemplate.New("").Funcs(xmlFuncMap)
	xmlTmpl, err = parseXMLFromLayers(xmlTmpl, layers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML templates: %w", err)
	}

	return &Renderer{
		blog:         b,
		singleTmpl:   singleTmpl,
		listTmpl:     listTmpl,
		xmlTemplates: xmlTmpl,
	}, nil
}

func buildTemplateLayers(b *Blog) []templateLayer {
	var layers []templateLayer

	layoutsDir := filepath.Join(b.root, "layouts")
	if info, err := os.Stat(layoutsDir); err == nil && info.IsDir() {
		layers = append(layers, templateLayer{
			name: "user",
			fsys: os.DirFS(layoutsDir),
			base: "",
		})
	}

	theme := b.Config.Theme
	if theme == "" {
		theme = "suspended"
	}

	themeLayoutsDir := filepath.Join(b.root, "themes", theme, "layouts")
	if info, err := os.Stat(themeLayoutsDir); err == nil && info.IsDir() {
		layers = append(layers, templateLayer{
			name: "theme:" + theme,
			fsys: os.DirFS(themeLayoutsDir),
			base: "",
		})
	}

	embeddedBase := fmt.Sprintf("themes/%s/layouts", theme)
	if sub, err := fs.Sub(embeddedThemes, embeddedBase); err == nil {
		layers = append(layers, templateLayer{
			name: "embedded:" + theme,
			fsys: sub,
			base: "",
		})
	}

	if theme != "suspended" {
		embeddedBase := "themes/suspended/layouts"
		if sub, err := fs.Sub(embeddedThemes, embeddedBase); err == nil {
			layers = append(layers, templateLayer{
				name: "embedded:suspended",
				fsys: sub,
				base: "",
			})
		}
	}

	return layers
}

func parsePartialsFromLayers(t *template.Template, layers []templateLayer) (*template.Template, error) {
	partialFiles := []string{
		"partials/head.html",
		"partials/header.html",
		"partials/footer.html",
		"partials/seo.html",
		"partials/styles.html",
		"partials/scripts.html",
	}

	parsed := make(map[string]bool)

	for _, file := range partialFiles {
		for _, layer := range layers {
			if parsed[file] {
				break
			}
			data, err := fs.ReadFile(layer.fsys, file)
			if err != nil {
				continue
			}
			t, err = t.Parse(string(data))
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s from %s: %w", file, layer.name, err)
			}
			parsed[file] = true
			break
		}
	}

	return t, nil
}

func parseFileFromLayers(t *template.Template, layers []templateLayer, file string) (*template.Template, error) {
	for _, layer := range layers {
		data, err := fs.ReadFile(layer.fsys, file)
		if err != nil {
			continue
		}
		return t.Parse(string(data))
	}
	return nil, fmt.Errorf("template %s not found in any layer", file)
}

func parseXMLFromLayers(t *texttemplate.Template, layers []templateLayer) (*texttemplate.Template, error) {
	xmlFiles := []string{"rss.xml", "sitemap.xml"}
	parsed := make(map[string]bool)

	for _, file := range xmlFiles {
		for _, layer := range layers {
			if parsed[file] {
				break
			}
			data, err := fs.ReadFile(layer.fsys, file)
			if err != nil {
				continue
			}
			t, err = t.New(file).Parse(string(data))
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s from %s: %w", file, layer.name, err)
			}
			parsed[file] = true
			break
		}
	}

	return t, nil
}

func IsBuiltinTheme(name string) bool {
	for _, t := range builtinThemes {
		if t == name {
			return true
		}
	}
	return false
}

func ListBuiltinThemes() []string {
	return append([]string{}, builtinThemes...)
}

func CopyBuiltinTheme(name, destDir string) error {
	if !IsBuiltinTheme(name) {
		return fmt.Errorf("unknown built-in theme: %s", name)
	}

	themeBase := fmt.Sprintf("themes/%s", name)

	return fs.WalkDir(embeddedThemes, themeBase, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(themeBase, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := fs.ReadFile(embeddedThemes, path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0644)
	})
}

type ThemeInfo struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Version     string `yaml:"version"`
	Repository  string `yaml:"repository"`
	Preview     string `yaml:"preview"`
	Path        string `yaml:"-"`
	IsBuiltin   bool   `yaml:"-"`
}

func LoadThemeInfo(themeDir string) (ThemeInfo, error) {
	var info ThemeInfo
	info.Path = themeDir

	themeYaml := filepath.Join(themeDir, "theme.yaml")
	data, err := os.ReadFile(themeYaml)
	if err != nil {
		info.Name = filepath.Base(themeDir)
		return info, nil
	}

	if err := parseYAML(data, &info); err != nil {
		return info, fmt.Errorf("failed to parse theme.yaml: %w", err)
	}

	return info, nil
}

func parseYAML(data []byte, v interface{}) error {
	lines := strings.Split(string(data), "\n")
	info := v.(*ThemeInfo)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")

		switch key {
		case "name":
			info.Name = value
		case "description":
			info.Description = value
		case "author":
			info.Author = value
		case "version":
			info.Version = value
		case "repository":
			info.Repository = value
		case "preview":
			info.Preview = value
		}
	}
	return nil
}

type PageData struct {
	Config       Config
	Article      *Article
	Articles     []Article
	HighlightCSS template.CSS
	CurrentYear  int
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
