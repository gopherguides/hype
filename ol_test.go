package hype

import "testing"

func Test_OL_MarshalJSON(t *testing.T) {
	t.Parallel()

	ol := &OL{
		Element: NewEl("ol", nil),
	}

	ol.Nodes = append(ol.Nodes, Text("This is an ordered list"))

	testJSON(t, "ol", ol)
}
