package hype

import "testing"

func Test_InlineCode_MarshalJSON(t *testing.T) {
	t.Parallel()

	il := &InlineCode{
		Element: NewEl("code", nil),
	}

	il.Nodes = append(il.Nodes, Text("var x = 1"))

	testJSON(t, "inline_code", il)

}
