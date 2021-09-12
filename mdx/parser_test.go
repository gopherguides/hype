package mdx

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var testdata = os.DirFS("testdata")

func Test_Parser_ElementNode_Custom(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	src, err := fs.ReadFile(testdata, "assignment.md")
	r.NoError(err)

	// act := blackfriday.Run(src, blackfriday.WithRenderer(&Renderer{}))

	p := New()
	act, err := p.Parse(src)
	r.NoError(err)

	exp := "<page>\n<h1>Assignment 42</h1>\n\n<assignment number=\"42\">\n\n<p>Instructions!</p>\n\n</assignment>\n</page>\n"

	// fmt.Println(string(act))
	r.Equal(exp, string(act))
}

func Test_Parser_Parse(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	src, err := fs.ReadFile(testdata, "module.md")
	r.NoError(err)

	p := New()

	out, err := p.Parse(src)
	r.NoError(err)

	act := string(out)
	act = strings.TrimSpace(act)

	exp := `<page>
<h1>Week 1 - Getting Started with Go - TODO</h1>

<p>overview</p>
</page>
<include src="basics.md"></include>
<include src="strings.md"></include>
<include src="utf8.md"></include>
<include src="numbers.md"></include>
<include src="booleans.md"></include>
<include src="variables.md"></include>
<include src="constants.md"></include>
<page>
<h1>Exercise (Due Wednesday)</h1>

<p>Locally create a new folder named, &ldquo;gopherguides-intro-to-go&rdquo;. Inside of the folder initialize a new git repository. Next, create a new folder named <code>gopherguides-intro-to-go/week01</code>. Inside of the <code>week01</code> folder initialize a new Go module named &ldquo;github.com/YOUR-USERNAME/gopherguides-intro-to-go/week01&rdquo;.</p>

<pre><code>$ go mod init github.com/YOUR-USERNAME/gopherguides-intro-to-go/week01
</code></pre>

<p>Commit the <code>go.mod</code> to the git repository.  Next, on <a href="https://github.com/">GitHub.com</a> create a new public repository named &ldquo;gopherguides-intro-to-go&rdquo; under your account and upload your local repository following the instructions on GitHub.</p>

<pre><code class="language-text">.
└── gopherguides-intro-to-go
    └── week01
        └── go.mod
</code></pre>
</page>
<page>
<h1>Assignment 1 (Due Sunday)</h1>

<h2>1.1</h2>

<p>Write a &ldquo;Hello, World&rdquo; style Go program using the <code>main</code> package. Your file should be named <code>main.go</code>. This program <strong>must</strong> compile and print &ldquo;Hello, World!&rdquo;, with a new line after it, to the console window when run. Publish this code to your repository under your <code>week01</code> folder you created earlier this week. Next, create a branch in your local project called, <code>assignment01</code>. Using <a href="pkg.go.dev">pkg.go.dev</a> research the <code>fmt</code> package. Use the <code>fmt</code> package to print &ldquo;Printing, TODO!&rdquo;, replacing &ldquo;TODO&rdquo; with the proper printing verb to properly print the following types: <code>string (&quot;Go&quot;), int (42), bool (true)</code>. Use <code>go vet</code> to confirm you are using the correct verb. Finally, open a PR to merge your new changes into your <code>main</code> branch. This PR should contain a paragraph or two explaining the changes and how they were implemented. Submit the link to the PR to be reviewed.</p>

<h2>1.2</h2>

<p>Write a short essay describing your history in technology and how you feel that Go fits into your plans for your future. Additionally, write a short essay discussing any surprises you found when researching the <code>fmt</code> package. Include how does printing in Go differ from other languages you may have used before. Please be specific and cite examples. (500 words minimum)</p>
</page>`

	// fmt.Println(act)
	r.Equal(exp, act)
}
