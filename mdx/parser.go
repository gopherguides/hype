package mdx

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/russross/blackfriday/v2"
)

const extensions = blackfriday.CommonExtensions | blackfriday.LaxHTMLBlocks

// Parser for parsing Markdown
type Parser struct {
	*sync.RWMutex
	cur   int
	lines []string
}

func (p *Parser) render(src []byte) []byte {
	act := blackfriday.Run(src, blackfriday.WithExtensions(extensions))
	return act
}

func (p *Parser) start(w io.Writer) {
	fmt.Fprintln(w, "<page>")
}

func (p *Parser) end(w io.Writer) {
	fmt.Fprintln(w, "</page>")
}

// NewParser returns a new Parser.
func New() *Parser {
	p := &Parser{
		RWMutex: &sync.RWMutex{},
	}

	return p
}

func (p *Parser) parse(lines []string) ([]byte, error) {
	bb := &bytes.Buffer{}

	var chunk []string
	var after string

	var ind int
	for _, line := range lines {
		ind++
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
		b, err := p.parse(lines[ind:])
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

// Parse parses the Markdown and returns the HTML.
func (p *Parser) Parse(src []byte) ([]byte, error) {
	p.Lock()
	p.lines = strings.Split(string(src), "\n")
	p.Unlock()

	return p.parse(p.lines)
}

// next returns the next line in the Markdown.
func (p *Parser) next() (string, bool) {
	p.Lock()
	ll := len(p.lines)
	cur := p.cur
	p.cur = p.cur + 1
	p.Unlock()

	if ll == 0 || cur >= ll {
		return "", false
	}

	// p.cur = p.cur + 1

	return p.lines[cur], true
}
