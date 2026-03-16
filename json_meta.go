package hype

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type DocumentMeta struct {
	Title              string            `json:"title"`
	Slug               string            `json:"slug"`
	Metadata           map[string]string `json:"metadata"`
	Headings           []HeadingMeta     `json:"headings"`
	CodeSnippets       []CodeSnippetMeta `json:"code_snippets"`
	Includes           []string          `json:"includes"`
	Images             []string          `json:"images"`
	WordCount          int               `json:"word_count"`
	ReadingTimeMinutes int               `json:"reading_time_minutes"`
}

type HeadingMeta struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
	ID    string `json:"id"`
}

type CodeSnippetMeta struct {
	Language string  `json:"language"`
	Src      *string `json:"src"`
	Snippet  *string `json:"snippet"`
}

type DocumentTOC struct {
	Title string     `json:"title"`
	TOC   []TOCEntry `json:"toc"`
}

type TOCEntry struct {
	Level    int        `json:"level"`
	Text     string     `json:"text"`
	ID       string     `json:"id"`
	Children []TOCEntry `json:"children"`
}

func ExtractMeta(doc *Document) *DocumentMeta {
	if doc == nil {
		return &DocumentMeta{}
	}

	meta := &DocumentMeta{
		Title:        doc.Title,
		Metadata:     make(map[string]string),
		Headings:     []HeadingMeta{},
		CodeSnippets: []CodeSnippetMeta{},
		Includes:     []string{},
		Images:       []string{},
	}

	mds := ByType[*Metadata](doc.Nodes)
	for _, md := range mds {
		for _, k := range md.Map.Keys() {
			v, _ := md.Map.Get(k)
			meta.Metadata[k] = v
		}
	}

	if doc.Parser != nil {
		for _, k := range doc.Parser.Vars.Keys() {
			if _, exists := meta.Metadata[k]; !exists {
				if v, ok := doc.Parser.Vars.Get(k); ok {
					meta.Metadata[k] = fmt.Sprintf("%v", v)
				}
			}
		}
	}

	if slug, ok := meta.Metadata["slug"]; ok {
		meta.Slug = slug
	} else {
		meta.Slug = slugify(meta.Title)
	}

	headings := ByType[*Heading](doc.Nodes)
	for _, h := range headings {
		text, id := extractHeadingTextAndID(h)
		meta.Headings = append(meta.Headings, HeadingMeta{
			Level: h.Level(),
			Text:  text,
			ID:    id,
		})
	}

	srcs := ByType[*SourceCode](doc.Nodes)
	for _, sc := range srcs {
		cs := CodeSnippetMeta{
			Language: sc.Lang,
		}
		if sc.Src != "" {
			s := sc.Src
			cs.Src = &s
		}
		if sc.Snippet.Name != "" {
			n := sc.Snippet.Name
			cs.Snippet = &n
		}
		meta.CodeSnippets = append(meta.CodeSnippets, cs)
	}

	fcs := ByType[*FencedCode](doc.Nodes)
	for _, fc := range fcs {
		cs := CodeSnippetMeta{
			Language: fc.Lang(),
		}
		meta.CodeSnippets = append(meta.CodeSnippets, cs)
	}

	incs := ByType[*Include](doc.Nodes)
	for _, inc := range incs {
		if src, ok := inc.Get("src"); ok {
			meta.Includes = append(meta.Includes, src)
		}
	}

	imgs := ByType[*Image](doc.Nodes)
	for _, img := range imgs {
		if src, ok := img.Get("src"); ok {
			meta.Images = append(meta.Images, src)
		}
	}

	meta.WordCount = countWords(doc)
	if meta.WordCount > 0 {
		meta.ReadingTimeMinutes = int(math.Ceil(float64(meta.WordCount) / 250.0))
	}

	return meta
}

func ExtractTOC(doc *Document) *DocumentTOC {
	if doc == nil {
		return &DocumentTOC{}
	}

	toc := &DocumentTOC{
		Title: doc.Title,
		TOC:   []TOCEntry{},
	}

	headings := ByType[*Heading](doc.Nodes)
	if len(headings) == 0 {
		return toc
	}

	var entries []HeadingMeta
	for _, h := range headings {
		text, id := extractHeadingTextAndID(h)
		entries = append(entries, HeadingMeta{
			Level: h.Level(),
			Text:  text,
			ID:    id,
		})
	}

	toc.TOC = buildTOCTree(entries, 0)
	return toc
}

var tocLevelRe = regexp.MustCompile(`<toc-level>[^<]*</toc-level>\s*-\s*`)
var anchorRe = regexp.MustCompile(`<a\s+id="([^"]*)">\s*</a>`)

func extractHeadingTextAndID(h *Heading) (text string, id string) {
	raw := h.Children().String()

	if m := anchorRe.FindStringSubmatch(raw); len(m) > 1 {
		id = m[1]
	}

	clean := anchorRe.ReplaceAllString(raw, "")
	clean = tocLevelRe.ReplaceAllString(clean, "")
	clean = stripHTMLFromMD(clean)
	text = strings.TrimSpace(clean)

	if id == "" {
		if attrID, ok := h.Get("id"); ok {
			id = attrID
		} else {
			id = slugify(text)
		}
	}

	return text, id
}

func buildTOCTree(entries []HeadingMeta, start int) []TOCEntry {
	var result []TOCEntry
	i := start

	for i < len(entries) {
		entry := TOCEntry{
			Level:    entries[i].Level,
			Text:     entries[i].Text,
			ID:       entries[i].ID,
			Children: []TOCEntry{},
		}

		j := i + 1
		for j < len(entries) && entries[j].Level > entries[i].Level {
			j++
		}

		if j > i+1 {
			entry.Children = buildTOCTree(entries[i+1:j], 0)
		}

		result = append(result, entry)
		i = j
	}

	return result
}

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	return s
}

func countWords(doc *Document) int {
	md := doc.Nodes.MD()

	md = stripHTMLFromMD(md)

	words := strings.Fields(md)
	return len(words)
}

func stripHTMLFromMD(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	return result.String()
}
