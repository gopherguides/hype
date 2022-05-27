package hype

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	goVersion = func() string {
		return "go.test"
	}
}

type brokenReader struct{}

func (brokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("broken reader")
}

func compareOutputFile(t testing.TB, cab fs.FS, act string, expFile string) {
	t.Helper()

	r := require.New(t)

	b, err := fs.ReadFile(cab, expFile)
	r.NoError(err)

	exp := string(b)

	compareOutput(t, act, exp)
}

func compareOutput(t testing.TB, act string, exp string) {
	t.Helper()

	r := require.New(t)

	// fn := func(s string) string {

	// 	// rx, err := regexp.Compile(`[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}`)

	// 	// r.NoError(err)

	// 	// uuids := rx.FindAllString(s, -1)
	// 	// for i, u := range uuids {
	// 	// 	s = strings.Replace(s, u, fmt.Sprintf("uuid-%d", i), -1)
	// 	// }

	// 	return strings.TrimSpace(s)
	// }

	// act = fn(act)

	// exp = fn(exp)

	// fmt.Println(act)
	r.Equal(exp, act)

}
