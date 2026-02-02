package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FencedCode_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	code := &FencedCode{
		Element: NewEl("code", nil),
	}
	code.Nodes = append(code.Nodes, Text("var x = 1"))

	r.NoError(code.Set("language", "go"))

	testJSON(t, "fenced_code", code)

}

func Test_FencedCode_MD(t *testing.T) {
	t.Parallel()

	t.Run("simple content uses backticks", func(t *testing.T) {
		r := require.New(t)

		code := &FencedCode{
			Element: NewEl("code", nil),
		}
		code.Nodes = Nodes{Text("var x = 1")}
		r.NoError(code.Set("language", "go"))

		md := code.MD()
		r.Contains(md, "```go")
		r.Contains(md, "var x = 1")
	})

	t.Run("content with triple backticks uses tildes", func(t *testing.T) {
		r := require.New(t)

		// Simulates an indented code block showing mermaid syntax
		code := &FencedCode{
			Element: NewEl("code", nil),
		}
		code.Nodes = Nodes{Text("```mermaid\ngraph LR\n    A --> B\n```")}

		md := code.MD()
		// Should use tildes to avoid conflicts
		r.True(md[0:3] == "~~~", "should start with tildes, got: %s", md[0:10])
		r.Contains(md, "```mermaid")
		r.True(md[len(md)-3:] == "~~~", "should end with tildes")
	})
}
