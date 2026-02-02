package hype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"strings"
)

type InlineCode struct {
	*Element
}

func (code *InlineCode) MarshalJSON() ([]byte, error) {
	if code == nil {
		return nil, ErrIsNil("inline code")
	}

	code.RLock()
	defer code.RUnlock()

	m, err := code.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(code)

	return json.MarshalIndent(m, "", "  ")
}

func (code *InlineCode) String() string {
	if code == nil || code.Element == nil {
		return "<code></code>"
	}

	bb := &bytes.Buffer{}

	fmt.Fprint(bb, code.StartTag())

	body := code.Nodes.String()
	body = html.EscapeString(body)
	fmt.Fprint(bb, body)
	fmt.Fprint(bb, code.EndTag())

	return bb.String()
}

func (code *InlineCode) MD() string {
	if code == nil || code.Element == nil {
		return ""
	}

	content := code.Nodes.String()

	// Multi-line content should be rendered as a fenced code block
	if strings.Contains(content, "\n") {
		return fmt.Sprintf("```\n%s\n```", strings.TrimSpace(content))
	}

	// Content with backticks needs special handling
	if strings.Contains(content, "`") {
		maxTicks := countMaxConsecutiveBackticks(content)
		fence := strings.Repeat("`", maxTicks+1)
		// Add spaces when content starts/ends with backtick
		if strings.HasPrefix(content, "`") || strings.HasSuffix(content, "`") {
			return fmt.Sprintf("%s %s %s", fence, content, fence)
		}
		return fmt.Sprintf("%s%s%s", fence, content, fence)
	}

	return fmt.Sprintf("`%s`", content)
}

// countMaxConsecutiveBackticks counts the maximum number of consecutive
// backticks in a string. This is used to determine how many backticks
// are needed to fence inline code that contains backticks.
func countMaxConsecutiveBackticks(s string) int {
	max, current := 0, 0
	for _, r := range s {
		if r == '`' {
			current++
			if current > max {
				max = current
			}
		} else {
			current = 0
		}
	}
	return max
}

func NewInlineCode(el *Element) (*InlineCode, error) {
	if el == nil {
		return nil, ErrIsNil("element")
	}

	code := &InlineCode{
		Element: el,
	}

	return code, nil
}

func NewInlineCodeNodes(p *Parser, el *Element) (Nodes, error) {
	code, err := NewInlineCode(el)
	if err != nil {
		return nil, err
	}

	return Nodes{code}, nil
}
