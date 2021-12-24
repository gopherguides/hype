package hype

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Snippet
type Snippet struct {
	Content  string // The content of the snippet
	File     string // the file name of the snippet
	Language string // the language of the snippet
	Name     string // the name of the snippet
	Start    int    // the start line of the snippet
	End      int    // the end line of the snippet
}

func (snip Snippet) String() string {
	return snip.Content
}

// Snippets is a map of Snippet
type Snippets map[string]Snippet

// Snippets returns a map of Snippets from the given file.
func (p *Parser) Snippets(src string, b []byte) (Snippets, error) {
	p.RLock()
	defer p.RUnlock()
	return ParseSnippets(src, b, p.snippetRules)

}

// SnippetRule sets a Sprintf string for a file extension.
// Example: SnippetRule(".html", "<!-- %s -->")
func (p *Parser) SnippetRule(ext string, rule string) {
	p.Lock()
	defer p.Unlock()
	p.snippetRules[ext] = rule
}

// ParseSnippets will parse the given src and return a map of Snippets.
func ParseSnippets(path string, src []byte, rules map[string]string) (Snippets, error) {
	if rules == nil {
		rules = map[string]string{}
	}
	snips := Snippets{}

	ext := filepath.Ext(path)
	rule, ok := rules[ext]
	if !ok {
		rule = "// %s"
	}
	open := map[string]Snippet{}
	for i, line := range strings.Split(string(src), "\n") {
		sl := strings.TrimSpace(line)

		pre := fmt.Sprintf(rule, "snippet:(.+)")
		rx, err := regexp.Compile(pre)
		if err != nil {
			return nil, err
		}

		if names := rx.FindStringSubmatch(sl); len(names) > 1 {
			name := names[1]
			name = strings.TrimSpace(name)
			snip, ok := open[name]
			if ok {
				snip.End = i + 1
				snips[name] = snip
				delete(open, name)
			}
			if !ok {
				snip.File = path
				snip.Language = strings.TrimPrefix(ext, ".")
				snip.Name = name
				snip.Start = i + 1
				open[name] = snip
			}

			snips[name] = snip
			continue
		}

		for k, snip := range open {
			snip.Content = strings.Join([]string{snip.Content, line}, "\n")
			open[k] = snip
		}

	}

	return snips, nil
}
