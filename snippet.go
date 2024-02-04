package hype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/markbates/sweets"
)

type Snippet struct {
	Content string `json:"content,omitempty"` // The content of the snippet
	File    string `json:"file,omitempty"`    // the file name of the snippet
	Lang    string `json:"lang,omitempty"`    // the language of the snippet
	Name    string `json:"name,omitempty"`    // the name of the snippet
	Start   int    `json:"start,omitempty"`   // the start line of the snippet
	End     int    `json:"end,omitempty"`     // the end line of the snippet
}

func (snip Snippet) MarshalJSON() ([]byte, error) {
	m := struct {
		Content string `json:"content,omitempty"`
		File    string `json:"file,omitempty"`
		Lang    string `json:"lang,omitempty"`
		Name    string `json:"name,omitempty"`
		Start   int    `json:"start,omitempty"`
		End     int    `json:"end,omitempty"`
		Type    string `json:"type,omitempty"`
	}{
		Content: snip.Content,
		File:    snip.File,
		Lang:    snip.Lang,
		Name:    snip.Name,
		Start:   snip.Start,
		End:     snip.End,
	}

	if snip.Content != "" {
		m.Type = fmt.Sprintf("%T", snip)
	}

	return json.Marshal(m)
}

func (snip Snippet) String() string {
	return snip.Content
}

func (snip Snippet) MD() string {
	return html.UnescapeString(snip.Content)
}

func (snip Snippet) IsZero() bool {
	return snip.Content == ""
}

func (snip Snippet) Children() Nodes {
	return nil
}

type Snippets struct {
	snippets map[string]map[string]Snippet
	rules    map[string]string

	once sync.Once
	mu   sync.RWMutex
}

func (sm *Snippets) MarshalJSON() ([]byte, error) {
	if sm == nil {
		return nil, ErrIsNil("snippets")
	}

	m := struct {
		Rules    map[string]string             `json:"rules,omitempty"`
		Snippets map[string]map[string]Snippet `json:"snippets,omitempty"`
		Type     string                        `json:"type,omitempty"`
	}{
		Rules:    sm.rules,
		Snippets: sm.snippets,
	}

	if len(sm.rules) > 0 || len(sm.snippets) > 0 {
		m.Type = fmt.Sprintf("%T", sm)
	}

	return json.Marshal(m)
}

func (sm *Snippets) UnmarshalJSON(data []byte) error {
	m := struct {
		Rules    map[string]string             `json:"rules,omitempty"`
		Snippets map[string]map[string]Snippet `json:"snippets,omitempty"`
	}{}

	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.snippets = m.Snippets
	sm.rules = m.Rules

	return nil
}

func (sm *Snippets) init() {
	if sm == nil {
		return
	}

	sm.once.Do(func() {
		sm.mu.Lock()
		defer sm.mu.Unlock()

		if sm.snippets == nil {
			sm.snippets = map[string]map[string]Snippet{}
		}

		if sm.rules == nil {
			sm.rules = map[string]string{
				".go":   "// %s",
				".html": "<!-- %s -->",
				".md":   "<!-- %s -->",
				".rb":   "# %s",
			}
		}
	})
}

func (sm *Snippets) Add(ext string, rule string) {
	if sm == nil {
		return
	}

	sm.init()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.rules[ext] = rule
}

func (sm *Snippets) Get(name string) (map[string]Snippet, bool) {
	if sm == nil {
		return nil, false
	}

	sm.init()

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if snips, ok := sm.snippets[name]; ok {
		return snips, true
	}

	return nil, false
}

func (sm *Snippets) TrimComments(path string, src []byte) ([]byte, error) {
	if sm == nil {
		return nil, ErrIsNil("snippets")
	}

	sm.init()

	ext := filepath.Ext(path)

	sm.mu.RLock()
	rule, ok := sm.rules[ext]
	if !ok {
		rule = "// %s"
	}

	sm.mu.RUnlock()

	pre := fmt.Sprintf(rule, "snippet:(.+)")

	rx, err := regexp.Compile(pre)
	if err != nil {
		return nil, err
	}

	var lines [][]byte

	for _, line := range bytes.Split(src, []byte("\n")) {
		sl := bytes.TrimSpace(line)

		if rx.Match(sl) {
			continue
		}

		lines = append(lines, line)
	}

	return bytes.Join(lines, []byte("\n")), nil
}

// ParseSnippets will parse the given src and return a map of Snippets.
func (sm *Snippets) Parse(path string, src []byte) (map[string]Snippet, error) {
	if sm == nil {
		return nil, fmt.Errorf("snippets is nil")
	}

	sm.init()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if snips, ok := sm.snippets[path]; ok {
		return snips, nil
	}

	snips := map[string]Snippet{}

	ext := filepath.Ext(path)

	rule, ok := sm.rules[ext]
	if !ok {
		rule = "// %s"
	}

	pre := fmt.Sprintf(rule, "snippet:(.+)")
	rx, err := regexp.Compile(pre)
	if err != nil {
		return nil, err
	}

	open := map[string]Snippet{}

	for i, line := range strings.Split(string(src), "\n") {
		sl := strings.TrimSpace(line)

		if names := rx.FindStringSubmatch(sl); len(names) > 1 {
			name := names[1]
			name = strings.TrimSpace(name)

			snip, ok := open[name]
			if ok {
				snip.End = i + 1
				if _, ok := snips[name]; ok {
					return nil, fmt.Errorf("duplicate snippet: %s#%s", path, name)
				}
				snip.Content = sweets.TrimLeftSpace(snip.Content)
				snip.Content = strings.TrimSpace(snip.Content)
				snips[name] = snip
				delete(open, name)
			} else {
				snip.File = path
				snip.Lang = strings.TrimPrefix(ext, ".")
				snip.Name = name
				snip.Start = i + 1
				open[name] = snip
			}

			continue
		}

		for k, snip := range open {
			snip.Content = strings.Join([]string{snip.Content, line}, "\n")
			open[k] = snip
		}

	}

	if len(open) > 0 {
		keys := make([]string, 0, len(open))
		for k := range open {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		return nil, fmt.Errorf("unclosed snippet: %s: %q", path, keys)
	}

	sm.snippets[path] = snips

	return snips, nil
}
