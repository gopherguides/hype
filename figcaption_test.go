package hype

import "testing"

func Test_Figcaption_MarshalJSON(t *testing.T) {
	t.Parallel()

	fig := &Figcaption{
		Element: NewEl("figcaption", nil),
	}
	fig.Nodes = append(fig.Nodes, Text("This is a caption"))

	testJSON(t, "figcaption", fig)

}
