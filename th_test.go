package hype

import "testing"

func Test_TH_MarshalJSON(t *testing.T) {
	t.Parallel()

	th := &TH{
		Element: NewEl("th", nil),
	}

	testJSON(t, "th", th)
}
