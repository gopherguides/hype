package hype

import "testing"

func Test_Figure_MarshalJSON(t *testing.T) {
	t.Parallel()

	f := &Figure{
		Element:   NewEl("figure", nil),
		Pos:       1,
		SectionID: 2,
		style:     "figure",
	}
	f.Nodes = append(f.Nodes, Text("This is a figure"))

	testJSON(t, "figure", f)

}
