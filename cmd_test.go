package hype

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_NewCmd_Errors(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	ats := &Attributes{}
	r.NoError(ats.Set("exec", ""))

	table := []struct {
		name string
		el   *Element
		e    error
	}{
		{name: "nil element", e: ErrIsNil("element")},
		{name: "missing exec", e: ErrAttrNotFound("exec"), el: &Element{}},
		{name: "empty exec", e: ErrAttrEmpty("exec"), el: &Element{
			Attributes: ats,
		}},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			_, err := NewCmd(tt.el)
			r.Error(err)
			r.True(errors.Is(err, tt.e))
		})
	}
}

func Test_Cmd_Execute(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	c := &Cmd{
		Element: NewEl("cmd", nil),
		Args:    []string{"echo", "hello"},
	}
	r.NoError(c.Set("exec", "echo hello"))
	r.NoError(c.Set("hide-cmd", ""))

	ctx := context.Background()
	doc := &Document{
		Parser: NewParser(nil),
	}

	doc.Nodes = append(doc.Nodes, c)
	err := doc.Execute(ctx)
	r.NoError(err)

	r.NotNil(c.Result())

	act := c.String()
	act = strings.TrimSpace(act)

	exp := `<cmd exec="echo hello" hide-cmd=""><pre><code class="language-shell" language="shell">hello</code></pre></cmd>`

	// fmt.Println(act)
	r.Equal(exp, act)
}

func Test_Cmd_Execute_UnexpectedExit(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	c := &Cmd{
		Element: NewEl("cmd", nil),
		Args:    []string{"go", "run", "main.go"},
	}

	r.NoError(c.Set("src", "testdata/commands/bad-exit"))

	ctx := context.Background()
	doc := &Document{
		Parser: NewParser(nil),
	}

	err := c.Execute(ctx, doc)
	r.Error(err)

	c.ExpectedExit = 1

	err = c.Execute(ctx, doc)
	r.NoError(err)

}

func Test_NewCmd(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	el := &Element{
		Attributes: &Attributes{},
	}

	r.NoError(el.Set("environ", "foo=bar,baz=qux"))
	r.NoError(el.Set("exec", "echo hello"))
	r.NoError(el.Set("exit", "1"))
	r.NoError(el.Set("src", "testdata/commands/bad-exit"))
	r.NoError(el.Set("timeout", "10ms"))

	c, err := NewCmd(el)
	r.NoError(err)

	r.NotNil(c)

	r.Equal([]string{"echo", "hello"}, c.Args)
	r.Equal([]string{"foo=bar", "baz=qux"}, c.Env)
	r.Equal(1, c.ExpectedExit)
	r.Equal(time.Millisecond*10, c.Timeout)

}

func Test_Cmd_Execute_Timeout(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	c := &Cmd{
		Element: &Element{
			Attributes: &Attributes{},
		},
		Args:    []string{"go", "run", "main.go"},
		Timeout: time.Millisecond * 5,
	}

	r.NoError(c.Set("src", "testdata/commands/timeout"))

	ctx := context.Background()
	doc := &Document{}

	err := c.Execute(ctx, doc)
	r.Error(err)

}
