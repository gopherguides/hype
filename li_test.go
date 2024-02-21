package hype

import "testing"

func Test_LI_MarshalJSON(t *testing.T) {
	t.Parallel()

	li := &LI{
		Element: NewEl("li", nil),
		Type:    "ul",
	}

	li.Nodes = append(li.Nodes, Text("This is a list item"))

	testJSON(t, "li", li)

}
