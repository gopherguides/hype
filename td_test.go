package hype

import "testing"

func Test_TD_MarshalJSON(t *testing.T) {
	t.Parallel()

	td := &TD{
		Element: NewEl("td", nil),
	}
	td.Nodes = append(td.Nodes, Text("foo"))

	testJSON(t, "td", td)
}
