package hype

import "testing"

func Test_Paragraph_MarshalJSON(t *testing.T) {
	t.Parallel()

	p := &Paragraph{
		Element: NewEl("p", nil),
	}

	p.Nodes = append(p.Nodes, Text("This is a paragraph"))

	testJSON(t, "p", p)
}
