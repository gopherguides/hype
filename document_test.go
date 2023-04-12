package hype

import (
	"context"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Document_Execute(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	mod := `# Page 1

<foo></foo>

<include src="second/second.md"></include>`

	second := `# Second Page

<foo></foo>`

	cab := fstest.MapFS{
		"module.md": &fstest.MapFile{
			Data: []byte(mod),
		},
		"second/second.md": &fstest.MapFile{
			Data: []byte(second),
		},
	}

	p := NewParser(cab)
	p.NodeParsers["foo"] = func(p *Parser, el *Element) (Nodes, error) {
		x := executeNode{
			Element: el,
		}
		x.ExecuteFn = func(ctx context.Context, d *Document) error {
			time.Sleep(time.Millisecond * 10)

			x.Lock()
			x.Nodes = append(x.Nodes, Text("baz"))
			x.Unlock()
			return nil
		}

		x.Nodes = append(x.Nodes, Text("bar"))

		return Nodes{x}, nil
	}

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	go func() {
		defer cancel()

		err = doc.Execute(ctx)
		r.NoError(err)
	}()

	<-ctx.Done()

	r.NotEqual(context.DeadlineExceeded, ctx.Err())

	act := doc.String()
	// fmt.Println(act)

	exp := `<html><head></head><body><page>
<h1>Page 1</h1>

<foo>barbaz</foo>
</page>
<page>
<h1>Second Page</h1>

<foo>barbaz</foo>
</page>

</body></html>`

	r.Equal(exp, act)

}

func Test_Document_MD(t *testing.T) {
	t.Skip("TODO: fix this test")
	t.Parallel()
	r := require.New(t)

	root := "testdata/doc/to_md"

	p := testParser(t, root)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "module.md")
	r.NoError(err)

	s := doc.MD()
	act := string(s)
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	b, err := fs.ReadFile(p.FS, "module.gold")
	r.NoError(err)

	exp := string(b)
	exp = strings.TrimSpace(exp)

	r.Equal(exp, act)
}
