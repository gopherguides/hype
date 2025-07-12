package hype

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gobuffalo/flect"
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

	// Check if we need to process any markdown content (like ``` code blocks)
	// If the figure content contains raw markdown that wasn't preprocessed,
	// we need to process it while preserving all elements
	body := f.Nodes.String()
	if strings.Contains(body, "```") {
		// Process the content through markdown preprocessing but preserve
		// the structure by re-parsing the result
		nodes, err := f.processMarkdownContent(p, body)
		if err != nil {
			return nil, f.WrapErr(err)
		}
		f.Nodes = nodes
	}

	return f, nil
}

// processMarkdownContent processes figure content that contains markdown syntax
// while preserving all elements including figcaption
func (f *Figure) processMarkdownContent(p *Parser, body string) (Nodes, error) {
	// Create a sub-parser with pages disabled to avoid wrapping figure content in <page> tags
	subParser, err := p.Sub(".")
	if err != nil {
		return nil, err
	}
	subParser.DisablePages = true

	// Run the content through the preprocessing pipeline to handle markdown syntax
	r := strings.NewReader(body)

	// Apply preprocessing (including markdown conversion) with pages disabled
	processedReader, err := subParser.PreParsers.PreParse(subParser, r)
	if err != nil {
		return nil, err
	}

	// Parse the preprocessed content as HTML
	htmlDoc, err := html.Parse(processedReader)
	if err != nil {
		return nil, err
	}

	// Extract the body content (preprocessing wraps content in <html><head></head><body>)
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
		return nil, fmt.Errorf("could not find body element after preprocessing")
	}

	// Parse each child of the body as a hype node
	var nodes Nodes
	for child := bodyElement.FirstChild; child != nil; child = child.NextSibling {
		// Skip pure whitespace text nodes for cleaner output
		if child.Type == html.TextNode && strings.TrimSpace(child.Data) == "" {
			continue
		}

		node, err := subParser.ParseHTMLNode(child, f)
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
