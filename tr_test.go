package hype

import "testing"

func Test_TR_MarshalJSON(t *testing.T) {
	t.Parallel()

	tr := &TR{
		Element: NewEl("tr", nil),
	}

	testJSON(t, "tr", tr)
}
