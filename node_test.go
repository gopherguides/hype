package hype

import (
	"testing"
	"testing/fstest"

	"github.com/markbates/fsx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func Test_Parser_ParseNode(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		node *html.Node
		err  bool
	}{
		{name: "comment", node: &html.Node{Type: html.CommentNode}},
		{name: "doctype", node: &html.Node{Type: html.DoctypeNode}},
		{name: "document", node: &html.Node{Type: html.DocumentNode}},
		{name: "element", node: &html.Node{Type: html.ElementNode}},
		{name: "error", node: &html.Node{Type: html.ErrorNode}, err: true},
		{name: "nil", node: nil, err: true},
		{name: "text", node: &html.Node{Type: html.TextNode}},
	}

	p := testParser(t, fstest.MapFS{})

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			tag, err := p.ParseNode(tt.node)
			if tt.err {
				r.Error(err)
				return
			}
			r.NoError(err)
			r.NotNil(tag)
		})
	}
}

func Test_Node_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab, err := fsx.DirFS("testdata")
	r.NoError(err)

	name := `html5.html`

	p := testParser(t, cab)

	doc, err := p.ParseFile(name)
	r.NoError(err)
	r.NotNil(doc)

	exp, err := p.ReadFile(name)
	r.NoError(err)
	r.NotEmpty(exp)

	act := doc.String()
	// fmt.Println(act)
	r.Equal(node_string_exp, act)
}

const (
	node_string_exp = `<!doctype html5>
<html lang="en"><head>
    <meta charset="utf-8" content="utf-8" property="charset" />
    <meta content="width=device-width, initial-scale=1" name="viewport" />

    <title>A Basic HTML5 Template</title>
    <meta content="A simple HTML5 Template for new projects." name="description" />
    <meta content="SitePoint" name="author" />

    <meta content="A Basic HTML5 Template" property="og:title" />
    <meta content="website" property="og:type" />
    <meta content="https://www.sitepoint.com/a-basic-html5-template/" property="og:url" />
    <meta content="A simple HTML5 Template for new projects." property="og:description" />
    <meta content="image.png" property="og:image" />

    <link href="/favicon.ico" rel="icon" />
    <link href="/favicon.svg" rel="icon" type="image/svg+xml" />
    <link href="/apple-touch-icon.png" rel="apple-touch-icon" />

    <link href="css/styles.css?v=1.0" rel="stylesheet" />

</head>

<body>


    <!--  your content here...  -->
    <script src="js/scripts.js"></script>

    <div class="text">
        <img src="assets/foo.png" width="100%" />
        <p>Hello World</p>
    </div>

    <p><pre class="code-block"><code class="language-go" language="go" snippet="main" src="src/main.go">func main() {</code></pre></p>
    <p><pre class="code-block"><code class="language-html" language="html" snippet="main" src="src/snip.html">&lt;p&gt;Hello World&lt;/p&gt;</code></pre></p>
    <p><pre class="code-block"><code class="language-txt" language="txt" snippet="main" src="src/snip.txt">Line 2</code></pre></p>




</body>
</html>`
)
