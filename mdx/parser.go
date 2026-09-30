package mdx

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/russross/blackfriday/v2"
)

const extensions = blackfriday.CommonExtensions | blackfriday.LaxHTMLBlocks

// Parser for parsing Markdown
type Parser struct {
	*sync.RWMutex
	DisablePages bool
	lines        []string
}

func (p *Parser) render(src []byte) []byte {
	act := blackfriday.Run(src, blackfriday.WithExtensions(extensions))
	return act
}

func (p *Parser) start(w io.Writer) {
	if !p.DisablePages {
		fmt.Fprintln(w, "<page>")
	}
}

func (p *Parser) end(w io.Writer) {
	if !p.DisablePages {
		fmt.Fprintln(w, "</page>")
	}
}

// NewParser returns a new Parser.
func New() *Parser {
	p := &Parser{
		RWMutex: &sync.RWMutex{},
	}

	return p
}

func (p *Parser) parse(lines []string, fenced []bool) ([]byte, error) {
	bb := &bytes.Buffer{}

	var chunk []string
	var after string

	var ind int
	for i, line := range lines {
		ind++
		if fenced[i] {
			chunk = append(chunk, line)
			continue
		}
		if strings.HasPrefix(line, "<include") {
			after = line
			break
		}
		if strings.HasPrefix(line, "---") {
			break
		}

		chunk = append(chunk, line)
	}

	in := []byte(strings.Join(chunk, "\n"))
	in = bytes.TrimSpace(in)

	if len(in) > 0 {
		p.start(bb)
		bb.Write(p.render(in))
		p.end(bb)
	}

	if len(after) > 0 {
		fmt.Fprintln(bb, after)
	}

	if ind < len(lines) {
		b, err := p.parse(lines[ind:], fenced[ind:])
		if err != nil {
			return nil, err
		}
		bb.Write(b)
	}

	rx, err := regexp.Compile("<p>(</?.+>)</p>")
	if err != nil {
		return nil, err
	}

	lines = []string{}
	for _, line := range strings.Split(bb.String(), "\n") {
		if m := rx.FindStringSubmatch(line); len(m) > 1 {
			lines = append(lines, m[1])
			continue
		}
		lines = append(lines, line)
	}

	act := []byte(strings.Join(lines, "\n"))
	return act, nil
}

func isBreak(line string) bool {
	return strings.HasPrefix(line, "<include") || strings.HasPrefix(line, "---")
}

var breakRx = regexp.MustCompile(`hypemdxbreak(\d+)x`)

func fencedLines(lines []string) []bool {
	marked := make([]string, len(lines))
	for i, line := range lines {
		if isBreak(line) {
			marked[i] = fmt.Sprintf("hypemdxbreak%dx", i)
			continue
		}
		marked[i] = line
	}

	fenced := make([]bool, len(lines))

	md := blackfriday.New(blackfriday.WithExtensions(extensions))
	doc := md.Parse([]byte(strings.Join(marked, "\n")))
	doc.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if n.Type != blackfriday.CodeBlock || !n.IsFenced {
			return blackfriday.GoToNext
		}
		for _, m := range breakRx.FindAllSubmatch(n.Literal, -1) {
			i, err := strconv.Atoi(string(m[1]))
			if err == nil && i < len(lines) && isBreak(lines[i]) {
				fenced[i] = true
			}
		}
		return blackfriday.GoToNext
	})

	return fenced
}

// Parse parses the Markdown and returns the HTML.
func (p *Parser) Parse(src []byte) ([]byte, error) {
	p.Lock()
	p.lines = strings.Split(string(src), "\n")
	p.Unlock()

	return p.parse(p.lines, fencedLines(p.lines))
}
