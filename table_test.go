package hype

import (
	"strings"
	"testing"

	"github.com/markbates/table"
	"github.com/stretchr/testify/require"
)

func Test_Table_Data(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	root := "testdata/table/data"
	p := testParser(t, root)

	doc, err := p.ParseFile("hype.md")
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

func Test_Table_MD_RowBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rows  []string
	}{
		{
			name: "inline code and empty cells",
			input: "| `Tag` | Example | Alias |\n" +
				"| --- | --- | --- |\n" +
				"| `<include>` | `<include src=\"other.md\">` | |\n" +
				"| `<cmd>` | `a  b` | `-x` |\n",
			rows: []string{
				"| `Tag` | Example | Alias |",
				"| ----- | ------- | ----- |",
				"| `<include>` | `<include src=\"other.md\">` |  |",
				"| `<cmd>` | `a  b` | `-x` |",
			},
		},
		{
			name: "multiline HTML cells",
			input: "<table><thead><tr><th>First\nLast</th><th>Example</th></tr></thead>" +
				"<tbody><tr><td>one\r\ntwo\rthree\nfour</td><td></td></tr>" +
				"<tr><td><code>Tag</code></td><td><code>&lt;img src=&quot;image.png&quot;&gt;</code></td></tr>" +
				"</tbody></table>",
			rows: []string{
				"| First Last | Example |",
				"| ---------- | ------- |",
				"| one two three four |  |",
				"| `Tag` | `<img src=\"image.png\">` |",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			p := testParser(t, "")
			doc, err := p.Parse(strings.NewReader(tt.input))
			r.NoError(err)

			md := strings.TrimSpace(doc.MD())
			r.Equal(tt.rows, strings.Split(md, "\n"))
			r.NotContains(md, "\r")
		})
	}
}

func Test_Table_No_THEAD(t *testing.T) {
	t.Parallel()
	t.Skip()
	root := "testdata/table/headless"

	testModule(t, root)
}

func Test_Table_MarshalJSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/table/data")

	doc, err := p.ParseFile("hype.md")
	r.NoError(err)

	tables := ByType[*Table](doc.Children())
	r.Len(tables, 1)

	tab := tables[0]
	r.NotNil(tab)

	testJSON(t, "table", tab)
}
