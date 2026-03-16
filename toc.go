package hype

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type ToC struct {
	*Element
	Depth   int
	Root    bool
	mdNodes string
}

func (toc *ToC) MarshalJSON() ([]byte, error) {
	if toc == nil {
		return nil, ErrIsNil("toc")
	}

	m, err := toc.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(toc)
	m["depth"] = toc.Depth
	m["root"] = toc.Root

	return json.MarshalIndent(m, "", "  ")
}

func (toc *ToC) MD() string {
	if toc == nil {
		return ""
	}

	if toc.mdNodes != "" {
		return toc.mdNodes
	}

	bb := &bytes.Buffer{}
	bb.WriteString(toc.StartTag())
	bb.WriteString("\n")
	bb.WriteString(toc.Nodes.MD())
	bb.WriteString(toc.EndTag())
	return bb.String()
}

func (toc *ToC) PostExecute(ctx context.Context, doc *Document, err error) error {
	if err != nil {
		return nil
	}

	if err := toc.validate(); err != nil {
		return err
	}

	if doc == nil {
		return toc.WrapErr(fmt.Errorf("document is nil"))
	}

	headings := ByType[*Heading](doc.Children())

	seen := map[string]int{}
	filtered := filterHeadings(headings, toc.Depth, toc.Root)

	slugs := make([]string, len(filtered))
	for i, h := range filtered {
		if existing, ok := h.Get("id"); ok && existing != "" {
			slugs[i] = existing
			seen[existing] = 1
		} else {
			text := h.Children().String()
			slugs[i] = UniqueSlug(text, seen)
		}
	}

	for i, h := range filtered {
		if _, ok := h.Get("id"); !ok {
			h.Set("id", slugs[i])
		}
	}

	nodes, err := GenerateToC(doc.Parser, filtered, slugs)
	if err != nil {
		return toc.WrapErr(err)
	}

	toc.Nodes = nodes

	toc.mdNodes = generateMDToC(filtered, slugs)

	return nil
}

func NewToCNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, fmt.Errorf("el is nil")
	}

	toc := &ToC{
		Element: el,
		Depth:   6,
		Root:    true,
	}

	if v, ok := el.Get("depth"); ok {
		d, err := strconv.Atoi(v)
		if err != nil {
			return nil, el.WrapErr(fmt.Errorf("invalid depth %q: %w", v, err))
		}
		toc.Depth = d
	}

	if v, ok := el.Get("root"); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, el.WrapErr(fmt.Errorf("invalid root %q: %w", v, err))
		}
		toc.Root = b
	}

	return Nodes{toc}, nil
}

func filterHeadings(headings []*Heading, depth int, root bool) []*Heading {
	var filtered []*Heading
	for _, h := range headings {
		if h.Level() > depth {
			continue
		}
		if !root && h.Level() == 1 {
			continue
		}
		filtered = append(filtered, h)
	}
	return filtered
}

func GenerateToC(p *Parser, headings []*Heading, slugs []string) (Nodes, error) {
	if len(headings) == 0 {
		return nil, nil
	}

	bb := &bytes.Buffer{}

	minLevel := headings[0].Level()
	for _, h := range headings[1:] {
		if h.Level() < minLevel {
			minLevel = h.Level()
		}
	}

	currentLevel := minLevel - 1
	bb.WriteString("<nav class=\"toc\">")

	for i, h := range headings {
		text := h.Children().String()
		slug := slugs[i]
		level := h.Level()

		if level > currentLevel {
			for level > currentLevel {
				indent := strings.Repeat("  ", currentLevel-minLevel+1)
				bb.WriteString("\n")
				bb.WriteString(indent)
				bb.WriteString("<ul>")
				currentLevel++
			}
		} else {
			bb.WriteString("</li>")
			for level < currentLevel {
				bb.WriteString("\n")
				bb.WriteString(strings.Repeat("  ", currentLevel-minLevel))
				bb.WriteString("</ul>")
				bb.WriteString("</li>")
				currentLevel--
			}
		}

		bb.WriteString("\n")
		bb.WriteString(strings.Repeat("  ", level-minLevel+1))
		fmt.Fprintf(bb, "<li><a href=\"#%s\">%s</a>", slug, text)
	}

	bb.WriteString("</li>")
	for currentLevel > minLevel {
		bb.WriteString("\n")
		bb.WriteString(strings.Repeat("  ", currentLevel-minLevel))
		bb.WriteString("</ul>")
		bb.WriteString("</li>")
		currentLevel--
	}
	bb.WriteString("\n</ul>\n</nav>")

	frag, err := p.ParseFragment(bb)
	if err != nil {
		return nil, err
	}

	for _, n := range frag.Children() {
		switch n.(type) {
		case Text:
		default:
			return Nodes{n}, nil
		}
	}

	return nil, fmt.Errorf("unable to parse toc")
}

func generateMDToC(headings []*Heading, slugs []string) string {
	bb := &bytes.Buffer{}

	minLevel := 7
	for _, h := range headings {
		if h.Level() < minLevel {
			minLevel = h.Level()
		}
	}

	for i, h := range headings {
		text := h.Children().String()
		slug := slugs[i]
		indent := strings.Repeat("  ", h.Level()-minLevel)
		fmt.Fprintf(bb, "%s- [%s](#%s)\n", indent, text, slug)
	}

	return bb.String()
}

func (toc *ToC) validate() error {
	if toc == nil {
		return toc.WrapErr(fmt.Errorf("toc is nil"))
	}

	if toc.Element == nil {
		return toc.WrapErr(fmt.Errorf("toc.Element is nil"))
	}

	return nil
}

