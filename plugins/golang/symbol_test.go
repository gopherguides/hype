package golang

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Symbols(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	const root = "testdata/symbols"
	cab := os.DirFS(root)

	p := testParser(t, cab, root)

	doc, err := p.ParseFile("symbols.md")
	r.NoError(err)
	r.NotNil(doc)

	act := doc.String()

	// fmt.Println(act)

	exp := `<html><head><meta charset="utf-8" /></head><body>
<page>

<h1>Tests for GoSrc</h1>

<p>Ask for a type</p>

<code class="language-go" language="go" src="src/foo" sym="User">type User struct {
	Name string
	Age  int
}

func (u User) String() string</code>

<p>Ask for a method on a type</p>

<code class="language-go" language="go" src="src/foo" sym="User.String">// String returns a string representation of the user.
func (u User) String() string {
	return fmt.Sprintf(&#34;%s (%d)&#34;, u.Name, u.Age)
}</code>

<p>Ask for a private constant</p>

<code class="language-go" language="go" src="src/foo" sym="orange">const (
	orange = iota
	apple
	pear
)</code>

</page><!--BREAK-->

<page>

<h1>Tests for Bar</h1>

<p>Ask for a method on a type</p>

<code class="language-go" language="go" src="bar/src/bar" sym="User.String">// String returns a string representation of the user.
func (u User) String() string {
	return fmt.Sprintf(&#34;%s (%d)&#34;, u.Name, u.Age)
}</code>

<p>ask for the full source for a type</p>

<code class="language-go" language="go" src="bar/src/bar" sym="-all User">type User struct {
	Name string
	Age  int
}

// String returns a string representation of the user.
func (u User) String() string {
	return fmt.Sprintf(&#34;%s (%d)&#34;, u.Name, u.Age)
}</code>

<p>show the main function</p>

<code class="language-go" language="go" src="bar/src/bar" sym="main">func main() {
	u := User{Name: &#34;jan&#34;, Age: 42}
	fmt.Println(u.String())
}</code>

</page><!--BREAK-->



</body>
</html>`

	r.Equal(exp, act)
}
