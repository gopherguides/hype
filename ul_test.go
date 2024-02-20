package hype

import "testing"

func Test_UL_MarshalJSON(t *testing.T) {
	t.Parallel()

	ul := &UL{
		Element: NewEl("ul", nil),
	}

	testJSON(t, "ul", ul)
}
