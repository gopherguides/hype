package hype

import (
	"os"
	"testing"

	"github.com/markbates/table"
	"github.com/stretchr/testify/require"
)

func Test_Table_Data(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/table/data"
	cab := os.DirFS(root)

	p := NewParser(cab)
	p.Root = root

	doc, err := p.ParseFile("module.md")
	r.NoError(err)

	// fmt.Println(doc.String())
	r.NotNil(doc)

	tables := ByType[*Table](doc.Children())
	r.Len(tables, 1)

	tab := tables[0]
	r.NotNil(tab)

	data, err := tab.Data()
	r.NoError(err)

	r.NotNil(data)

	cols, err := data.Columns()
	r.NoError(err)

	r.Len(cols, 2)
	r.Equal([]string{"Name", "Age"}, cols)

	rows, err := data.Rows()
	r.NoError(err)

	r.Len(rows, 3)

	r.Equal(table.Row{"Alice", "42"}, rows[0])
	r.Equal(table.Row{"Bob", "13"}, rows[1])
	r.Equal(table.Row{"Kurt", "27"}, rows[2])
}

func Test_Table_MD_in_MD(t *testing.T) {
	t.Parallel()

	root := "testdata/table/md_in_md"

	testModule(t, root)
}

func Test_Table_MD_in_HTML(t *testing.T) {
	t.Parallel()

	root := "testdata/table/md_in_html"

	testModule(t, root)
}

func Test_Table_No_THEAD(t *testing.T) {
	t.Parallel()
	t.Skip()
	root := "testdata/table/headless"

	testModule(t, root)
}
