package golang

import (
	"fmt"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Envelope_FilePath(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	full := Envelope{
		Body: "Hello, World",
		Cmd:  exec.Command("go", "run", "main.go"),
		Doc:  "context#Context",
		Root: "testdata/cmd",
		Exit: -1,
	}

	act, err := full.FilePath()
	r.NoError(err)

	exp := `TODO`
	r.Equal(exp, act)
}

func Test_Envelope_String(t *testing.T) {
	t.Parallel()
	t.Skip()
	v := runtime.Version()

	body := "Hello, World!"

	c := exec.Command("go", "run", "main.go")

	cmd := Envelope{
		Body: body,
		Cmd:  c,
	}

	doc := Envelope{
		Body: body,
		Doc:  "context#Context",
	}

	full := Envelope{
		Body: body,
		Cmd:  c,
		Doc:  "context#Context",
		Root: "src/foo",
		Exit: -1,
	}

	emptyExp := fmt.Sprintf("Envelope:\n\tgo version:\t%s", v)

	cmdExp := fmt.Sprintf("Envelope:\n\tgo version:\t%s\n\tcommand:\t$ go run main.go\n\n---\n\nHello, World!", v)

	docExp := fmt.Sprintf("Envelope:\n\tgo version:\t%s\n\tdocs:\t\t<a href=\"https://pkg.go.dev/context#Context\" target='_blank'>context#Context</a>\n\n---\n\nHello, World!", v)

	fullExp := fmt.Sprintf("Envelope:\n\tgo version:\t%s\n\tcommand:\t$ go run main.go\n\troot:\t\tsrc/foo\n\tdocs:\t\t<a href=\"https://pkg.go.dev/context#Context\" target='_blank'>context#Context</a>\n\n---\n\nHello, World!", v)

	table := []struct {
		name string
		env  Envelope
		exp  string
	}{
		{name: "cmd", env: cmd, exp: cmdExp},
		{name: "doc", env: doc, exp: docExp},
		{name: "empty", exp: emptyExp},
		{name: "full", env: full, exp: fullExp},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			_ = v

			e := tt.env

			act := e.String()
			fmt.Println(act)
			r.Equal(tt.exp, act)
		})
	}
}
