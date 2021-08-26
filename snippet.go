package hype

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type Snippet struct {
	Content  string
	File     string
	Language string
	Name     string
	Start    int
	End      int
}

func (snip Snippet) String() string {
	return snip.Content
}

type Snippets map[string]Snippet

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

func ParseSnippets(src string, b []byte, rules map[string]string) (Snippets, error) {
	if rules == nil {
		rules = map[string]string{}
	}
	snips := Snippets{}

	ext := filepath.Ext(src)
	rule, ok := rules[ext]
	if !ok {
		rule = "// %s"
	}
	open := map[string]Snippet{}
	for i, line := range strings.Split(string(b), "\n") {
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
				snip.File = src
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
