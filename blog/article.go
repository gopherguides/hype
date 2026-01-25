package blog

import (
	"context"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"math"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gopherguides/hype"
)

type Article struct {
	Title          string            `json:"title"`
	Slug           string            `json:"slug"`
	Published      time.Time         `json:"published"`
	Author         string            `json:"author"`
	Tags           []string          `json:"tags"`
	SEODescription string            `json:"seo_description"`
	OGImage        string            `json:"og_image"`
	Body           template.HTML     `json:"body"`
	Overview       template.HTML     `json:"overview"`
	Data           map[string]string `json:"data"`
	ReadingTime    int               `json:"reading_time"`
	File           string            `json:"file"`
	Dir            string            `json:"dir"`
}

func (a Article) IsPublished() bool {
	return time.Now().After(a.Published)
}

func (a Article) URL() string {
	return "/" + a.Slug + "/"
}

func (a Article) FormattedDate() string {
	return a.Published.Format("January 2, 2006")
}

func (a Article) ISODate() string {
	return a.Published.Format(time.RFC3339)
}

type ArticleParser struct {
	fsys        fs.FS
	highlighter *Highlighter
	baseURL     string
}

func NewArticleParser(fsys fs.FS, highlighter *Highlighter, baseURL string) *ArticleParser {
	return &ArticleParser{
		fsys:        fsys,
		highlighter: highlighter,
		baseURL:     baseURL,
	}
}

func (ap *ArticleParser) ParseArticle(ctx context.Context, dir string) (Article, error) {
	a := Article{
		Data: make(map[string]string),
		Dir:  dir,
	}

	subFS, err := fs.Sub(ap.fsys, dir)
	if err != nil {
		return a, fmt.Errorf("failed to access directory %s: %w", dir, err)
	}

	mdFile := "module.md"
	if _, err := fs.Stat(subFS, mdFile); err != nil {
		mdFile = "hype.md"
		if _, err := fs.Stat(subFS, mdFile); err != nil {
			entries, _ := fs.ReadDir(subFS, ".")
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".md") {
					mdFile = e.Name()
					break
				}
			}
		}
	}
	a.File = filepath.Join(dir, mdFile)

	p := hype.NewParser(subFS)
	doc, err := p.ParseFile(mdFile)
	if err != nil {
		return a, fmt.Errorf("failed to parse %s: %w", a.File, err)
	}

	if err := doc.Execute(ctx); err != nil {
		return a, fmt.Errorf("failed to execute %s: %w", a.File, err)
	}

	p.Vars.Range(func(key string, value any) bool {
		if s, ok := value.(string); ok {
			a.Data[key] = s
		}
		return true
	})

	slug := filepath.Base(dir)
	if v, ok := p.Vars.Get("slug"); ok {
		if s, ok := v.(string); ok {
			slug = s
		}
	}
	a.Slug = slug

	pbVal, ok := p.Vars.Get("published")
	if !ok {
		return a, fmt.Errorf("missing published date in %s", a.File)
	}
	pb, ok := pbVal.(string)
	if !ok {
		return a, fmt.Errorf("published date is not a string in %s", a.File)
	}

	t, err := time.Parse("01/02/2006", pb)
	if err != nil {
		return a, fmt.Errorf("invalid published date format in %s: %w", a.File, err)
	}
	a.Published = t

	if v, ok := p.Vars.Get("author"); ok {
		if s, ok := v.(string); ok {
			a.Author = s
		}
	}

	if v, ok := p.Vars.Get("seo_description"); ok {
		if s, ok := v.(string); ok {
			a.SEODescription = s
		}
	}

	if v, ok := p.Vars.Get("og_image"); ok {
		if s, ok := v.(string); ok {
			a.OGImage = s
		}
	}

	if v, ok := p.Vars.Get("tags"); ok {
		if s, ok := v.(string); ok {
			for _, tag := range strings.Split(s, ",") {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					a.Tags = append(a.Tags, tag)
				}
			}
		}
	}

	a.Title = doc.Title

	body, err := doc.Body()
	if err != nil {
		return a, fmt.Errorf("failed to get body from %s: %w", a.File, err)
	}

	bodyStr := body.String()
	bodyStr = stripDetailsBlocks(bodyStr)
	bodyStr = ap.processCodeBlocks(bodyStr)
	a.Body = template.HTML(bodyStr)

	a.ReadingTime = calculateReadingTime(bodyStr)

	pages, err := doc.Pages()
	if err == nil {
		for _, page := range pages {
			if v, ok := p.Vars.Get("overview"); ok {
				if s, ok := v.(string); ok && s == "true" {
					a.Overview = template.HTML(page.Children().String())
					break
				}
			}
		}
	}

	return a, nil
}

var codeBlockPattern = regexp.MustCompile(`<pre><code class="language-(\w+)"[^>]*>([\s\S]*?)</code></pre>`)

func (ap *ArticleParser) processCodeBlocks(htmlContent string) string {
	if ap.highlighter == nil {
		return htmlContent
	}

	return codeBlockPattern.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := codeBlockPattern.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		language := submatches[1]
		code := html.UnescapeString(submatches[2])

		highlighted, err := ap.highlighter.Highlight(code, language)
		if err != nil {
			return match
		}

		return highlighted
	})
}

var detailsPattern = regexp.MustCompile(`(?s)<details[^>]*>.*?</details>`)

func stripDetailsBlocks(html string) string {
	return detailsPattern.ReplaceAllString(html, "")
}

func calculateReadingTime(content string) int {
	plainText := stripHTML(content)
	words := len(strings.Fields(plainText))
	wordsPerMinute := 200.0
	minutes := float64(words) / wordsPerMinute
	return int(math.Ceil(minutes))
}

func stripHTML(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			result.WriteRune(r)
		}
	}
	return result.String()
}
