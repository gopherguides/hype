package hype

import "testing"

func Test_THead_MarshalJSON(t *testing.T) {
	t.Parallel()

	th := &THead{
		Element: NewEl("th", nil),
	}

	testJSON(t, "thead", th)

}
