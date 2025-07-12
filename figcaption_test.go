package hype

import (
	"context"
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

// Test that verifies the fix for figcaption elements being lost when mixed with <go> elements
func Test_Figcaption_With_Go_Elements(t *testing.T) {
	t.Parallel()

	// This structure was previously failing with "execute error: no figcaption"
	input := `<figure id="ticker" type="listing">
<go doc="time.Ticker"></go>
<figcaption>The <godoc>time#Ticker</godoc> function.</figcaption>
</figure>`

	p := testParser(t, "")

	doc, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error during parsing: %v", err)
	}

	// Execute the document - this should now work without the "no figcaption" error
	ctx := context.Background()
	err = doc.Execute(ctx)
	if err != nil {
		t.Fatalf("unexpected error during execution: %v", err)
	}

	// Verify that the figcaption is properly parsed and available
	figures := ByType[*Figure](doc.Nodes)
	if len(figures) == 0 {
		t.Fatal("no figures found")
	}

	fig := figures[0]
	figcaptions := ByType[*Figcaption](fig.Nodes)
	if len(figcaptions) == 0 {
		t.Fatal("no figcaption found in figure - the fix didn't work")
	}

	if len(figcaptions) != 1 {
		t.Fatalf("expected 1 figcaption, got %d", len(figcaptions))
	}
}
