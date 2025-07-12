package hype

import (
	"strings"
	"testing"
)

func Test_Figcaption_MarshalJSON(t *testing.T) {
	t.Parallel()

	fig := &Figcaption{
		Element: NewEl("figcaption", nil),
	}
	fig.Nodes = append(fig.Nodes, Text("This is a caption"))

	testJSON(t, "figcaption", fig)

}

func Test_Figcaption_WithGodoc(t *testing.T) {
	t.Parallel()

	input := `<figure id="ticker" type="listing">\n<go doc="time.Ticker"></go>\n<figcaption>The <godoc>time#Ticker</godoc> function.</figcaption>\n</figure>`

	p := NewParser(nil)
	nodes, err := p.ParseFragment(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseFragment error: %v", err)
	}

	figures := ByType[*Figure](nodes)
	if len(figures) == 0 {
		t.Fatalf("No Figure found in nodes")
	}
	fig := figures[0]

	caps := ByType[*Figcaption](fig.Nodes)
	if len(caps) == 0 {
		t.Fatalf("No Figcaption found in Figure nodes")
	}
	fc := caps[0]

	if got := fc.Nodes.String(); !strings.Contains(got, "Ticker") {
		t.Errorf("Figcaption content = %q, want to contain %q", got, "Ticker")
	}
}
