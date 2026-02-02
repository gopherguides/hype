package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_InlineCode_MarshalJSON(t *testing.T) {
	t.Parallel()

	il := &InlineCode{
		Element: NewEl("code", nil),
	}

	il.Nodes = append(il.Nodes, Text("var x = 1"))

	testJSON(t, "inline_code", il)

}

func Test_InlineCode_MD(t *testing.T) {
	t.Parallel()

	// Note: Multi-line content is now routed to FencedCode at parse time
	// InlineCode only handles single-line inline code
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple inline code",
			content:  "var x = 1",
			expected: "`var x = 1`",
		},
		{
			name:     "content with single backtick",
			content:  "use `fmt.Println`",
			expected: "`` use `fmt.Println` ``",
		},
		{
			name:     "content with double backticks",
			content:  "use ``code``",
			expected: "``` use ``code`` ```",
		},
		{
			name:     "content starting with backtick",
			content:  "`start",
			expected: "`` `start ``",
		},
		{
			name:     "content ending with backtick",
			content:  "end`",
			expected: "`` end` ``",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "``",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			code := &InlineCode{
				Element: NewEl("code", nil),
			}
			code.Nodes = Nodes{Text(tt.content)}

			act := code.MD()
			r.Equal(tt.expected, act)
		})
	}
}

func Test_InlineCode_MD_Nil(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	var code *InlineCode
	r.Equal("", code.MD())

	code = &InlineCode{}
	r.Equal("", code.MD())
}
