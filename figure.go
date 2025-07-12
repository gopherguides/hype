package hype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype/mdx"
	"golang.org/x/net/html"
)

type Figure struct {
	*Element

	Pos       int
	SectionID int

	style string
}

func (f *Figure) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, ErrIsNil("figure")
	}

	f.RLock()
	defer f.RUnlock()

	m, err := f.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(f)
	m["pos"] = f.Pos
	m["section_id"] = f.SectionID
	m["style"] = f.Style()

	return json.MarshalIndent(m, "", "  ")
}

func (f *Figure) MD() string {
	if f == nil {
		return ""
	}

	bb := &strings.Builder{}
	// bb.WriteString(f.StartTag())
	fmt.Fprintf(bb, "<a id=%q></a>\n", strings.TrimPrefix(f.Link(), "#"))
	fmt.Fprintln(bb)
	bb.WriteString(f.Nodes.MD())
	// bb.WriteString(f.EndTag())

	return bb.String()
}

// Style returns type of the figure.
// ex: "figure", "table", "listing", ...
func (f *Figure) Style() string {
	if f == nil || len(f.style) == 0 {
		return "figure"
	}

	return f.style
}

func (f *Figure) Name() string {
	if f == nil {
		return ""
	}

	style := flect.Titleize(f.Style())

	return fmt.Sprintf("%s %d.%d", style, f.SectionID, f.Pos)
}

func (f *Figure) Link() string {
	if f == nil {
		return ""
	}

	id, _ := f.Get("id")
	return fmt.Sprintf("#%s", id)
}

func NewFigure(p *Parser, el *Element) (*Figure, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	if p == nil {
		return nil, ErrIsNil("parser")
	}

	f := &Figure{
		Element: el,
	}

	if _, err := f.ValidAttr("id"); err != nil {
		return nil, err
	}

	f.style = "figure"
	style, ok := f.Get("type")
	style = strings.TrimSpace(style)

	if ok && len(style) > 0 {
		f.style = style
	}

	body := f.Nodes.String()
	body = strings.TrimSpace(body)

	if len(body) == 0 {
		return f, nil
	}

	// Create a sub-parser for processing the nodes
	p2, err := p.Sub(".")
	if err != nil {
		return nil, f.WrapErr(err)
	}

	// Parse figure content directly as HTML to avoid Markdown preprocessing issues
	// that can cause figcaption elements to be lost
	nodes, err := f.parseContentDirectly(p2, body)
	if err != nil {
		return nil, f.WrapErr(err)
	}

	f.Nodes = nodes
	return f, nil
}

// parseContentDirectly parses figure content with markdown preprocessing but avoids
// the paragraph extraction that can lose figcaption elements
func (f *Figure) parseContentDirectly(p *Parser, body string) (Nodes, error) {
	// Create a custom markdown preprocessor with DisablePages = true to avoid
	// wrapping figure content in <page> tags

	// Create a reader from the body content
	r := strings.NewReader(body)

	// Create a custom markdown preprocessor that doesn't add page tags
	md := mdx.New()
	md.DisablePages = true

	// Read the content
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Process with markdown (handles ``` code blocks, etc.)
	b = bytes.ReplaceAll(b, []byte("\\n"), []byte("  \n"))

	b, err = md.Parse(b)
	if err != nil {
		return nil, err
	}

	// Apply the same post-processing as the main markdown processor
	b = bytes.ReplaceAll(b, []byte("&rsquo;"), []byte("'"))
	b = bytes.ReplaceAll(b, []byte("&ldquo;"), []byte("\""))
	b = bytes.ReplaceAll(b, []byte("&rdquo;"), []byte("\""))

	// Parse the markdown-processed content as HTML
	htmlDoc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	// Find the body element and extract its children
	// The structure should be: html -> head, body -> content
	var bodyElement *html.Node
	var findBody func(*html.Node)
	findBody = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			bodyElement = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			findBody(child)
			if bodyElement != nil {
				return
			}
		}
	}
	findBody(htmlDoc)

	if bodyElement == nil {
		return nil, fmt.Errorf("could not find body element")
	}

	// Parse each child of the body as a node, skipping pure whitespace text nodes
	var nodes Nodes
	for child := bodyElement.FirstChild; child != nil; child = child.NextSibling {
		// Skip text nodes that contain only whitespace
		if child.Type == html.TextNode && strings.TrimSpace(child.Data) == "" {
			continue
		}

		node, err := p.ParseHTMLNode(child, f)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func NewFigureNodes(p *Parser, el *Element) (Nodes, error) {
	if p == nil {
		return nil, ErrIsNil("parser")
	}

	f, err := NewFigure(p, el)
	if err != nil {
		return nil, err
	}

	f.Lock()
	f.SectionID = p.Section
	f.Unlock()

	return Nodes{f}, nil
}
