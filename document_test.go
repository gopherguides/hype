package hype

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Document_Execute(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		r := require.New(t)
		p := testParser(t, "testdata/doc/execution/success")

		doc, err := p.ParseFile("hype.md")
		r.NoError(err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err = doc.Execute(ctx)
		r.NoError(err)

		act := doc.String()
		act = strings.TrimSpace(act)

		exp := "<html><head></head><body><page>\n<h1>Command</h1>\n\n<cmd exec=\"echo 'Hello World'\"><pre><code class=\"language-shell\" language=\"shell\">$ echo Hello World\n\nHello World</code></pre></cmd>\n</page>\n</body></html>"

		r.Equal(exp, act)
	})

	t.Run("failure", func(t *testing.T) {
		r := require.New(t)
		p := testParser(t, "testdata/doc/execution/failure")

		doc, err := p.ParseFile("hype.md")
		r.NoError(err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err = doc.Execute(ctx)
		r.Error(err)

		_, ok := err.(ExecuteError)
		r.True(ok, err)
		r.True(errors.Is(err, ExecuteError{}), err)
	})

}

func Test_Document_MD(t *testing.T) {
	t.Skip("TODO: fix this test")
	t.Parallel()
	r := require.New(t)

	root := "testdata/doc/to_md"

	p := testParser(t, root)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	s := doc.MD()
	act := string(s)
	act = strings.TrimSpace(act)

	// fmt.Println(act)

	b, err := fs.ReadFile(p.FS, "hype.gold")
	r.NoError(err)

	exp := string(b)
	exp = strings.TrimSpace(exp)

	r.Equal(exp, act)
}

func Test_Document_MarshalJSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/doc/snippets")
	p.DocIDGen = func() (string, error) {
		return "1", nil
	}

	ctx := context.Background()

	doc, err := p.ParseExecuteFile(ctx, "hype.md")
	r.NoError(err)

	testJSON(t, "document", doc)
}

func Test_Document_Pages(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/doc/pages")

	doc, err := p.ParseFile("hype.md")
	r.NoError(err)

	pages, err := doc.Pages()
	r.NoError(err)

	r.Len(pages, 3)
}

func Test_Document_Pages_NoPages(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	p := testParser(t, "testdata/doc/simple")

	doc, err := p.ParseFile("hype.md")
	r.NoError(err)

	pages, err := doc.Pages()
	r.NoError(err)

	r.Len(pages, 1)
}

// Test_Document_JSON_Determinism verifies that marshaling the same document to JSON
// produces identical output every time, ensuring predictable/deterministic behavior.
func Test_Document_JSON_Determinism(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		content string
		runs    int
	}{
		{
			name: "basic document",
			content: `<html><body><page>
<h1>Test Title</h1>
<p>This is a test paragraph.</p>
</page></body></html>`,
			runs: 10,
		},
		{
			name: "document with attributes",
			content: `<html><body><page>
<h1 id="title" class="header">Test Title</h1>
<p class="content" data-test="value">Paragraph with attributes.</p>
<div id="container" class="wrapper" data-foo="bar" data-baz="qux">Content</div>
</page></body></html>`,
			runs: 10,
		},
		{
			name: "document with code snippets",
			content: `<html><body><page>
<h1>Code Example</h1>
<pre><code class="language-go">
package main

func main() {
    println("Hello, World!")
}
</code></pre>
</page></body></html>`,
			runs: 10,
		},
		{
			name: "complex document",
			content: `<html><body><page>
<h1 id="main" class="title">Complex Document</h1>
<p class="intro">Introduction paragraph.</p>
<div id="content" class="section">
  <h2 class="subtitle">Section 1</h2>
  <p>Some content here.</p>
</div>
<div id="code" class="section">
  <h2>Code Section</h2>
  <pre><code class="language-go">func test() {}</code></pre>
</div>
</page></body></html>`,
			runs: 10,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			// Parse the document multiple times and marshal to JSON
			var jsonOutputs [][]byte
			for i := 0; i < tc.runs; i++ {
				cab := fstest.MapFS{
					"test.md": &fstest.MapFile{
						Data: []byte(tc.content),
					},
				}

				p := NewParser(cab)
				// Use a fixed ID generator to ensure deterministic IDs
				p.DocIDGen = func() (string, error) {
					return "test-doc-id", nil
				}

				doc, err := p.ParseFile("test.md")
				r.NoError(err)

				jsonData, err := json.Marshal(doc)
				r.NoError(err)
				jsonOutputs = append(jsonOutputs, jsonData)
			}

			// Verify all JSON outputs are identical
			firstOutput := jsonOutputs[0]
			for i := 1; i < len(jsonOutputs); i++ {
				if !bytes.Equal(firstOutput, jsonOutputs[i]) {
					t.Errorf("Run %d produced different JSON output:\nFirst:\n%s\n\nRun %d:\n%s\n",
						1, string(firstOutput), i+1, string(jsonOutputs[i]))
				}
			}
		})
	}
}

// Test_Document_JSON_Determinism_WithExecution tests that JSON output remains
// deterministic even after document execution.
func Test_Document_JSON_Determinism_WithExecution(t *testing.T) {
	t.Parallel()

	content := `<html><body><page>
<h1 id="title" class="header">Test With Execution</h1>
<p class="content">Content paragraph.</p>
</page></body></html>`

	runs := 10
	var jsonOutputs [][]byte

	for i := 0; i < runs; i++ {
		cab := fstest.MapFS{
			"test.md": &fstest.MapFile{
				Data: []byte(content),
			},
		}

		p := NewParser(cab)
		// Use fixed ID and time for determinism
		p.DocIDGen = func() (string, error) {
			return "test-doc-id", nil
		}
		p.NowFn = func() time.Time {
			return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		doc, err := p.ParseFile("test.md")
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		err = doc.Execute(ctx)
		cancel()
		require.NoError(t, err)

		jsonData, err := json.Marshal(doc)
		require.NoError(t, err)
		jsonOutputs = append(jsonOutputs, jsonData)
	}

	// Verify all JSON outputs are identical
	firstOutput := jsonOutputs[0]
	for i := 1; i < len(jsonOutputs); i++ {
		if !bytes.Equal(firstOutput, jsonOutputs[i]) {
			t.Errorf("Run %d produced different JSON output:\nFirst:\n%s\n\nRun %d:\n%s\n",
				1, string(firstOutput), i+1, string(jsonOutputs[i]))
		}
	}
}

// Test_Document_JSON_Determinism_RealWorld tests with a more realistic document
// that includes multiple features.
func Test_Document_JSON_Determinism_RealWorld(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	// Create a test filesystem with files
	cab := fstest.MapFS{
		"hype.md": &fstest.MapFile{
			Data: []byte(`<html><body><page>
<h1 id="main-title" class="title">Real World Test</h1>
<p class="intro">This is a real world test document.</p>
<div id="section1" class="section">
  <h2 class="subtitle">Section 1</h2>
  <p class="content">Content for section 1.</p>
</div>
<div id="section2" class="section">
  <h2 class="subtitle">Section 2</h2>
  <pre><code class="language-go">
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
</code></pre>
</div>
</page></body></html>`),
		},
	}

	runs := 10
	var jsonOutputs [][]byte

	for i := 0; i < runs; i++ {
		p := NewParser(cab)
		p.Root = "/test"
		// Use fixed ID for determinism
		p.DocIDGen = func() (string, error) {
			return "test-doc-id", nil
		}
		p.NowFn = func() time.Time {
			return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		// Add some vars
		p.Vars.Set("var1", "value1")
		p.Vars.Set("var2", "value2")

		doc, err := p.ParseFile("hype.md")
		r.NoError(err)

		jsonData, err := json.Marshal(doc)
		r.NoError(err)
		jsonOutputs = append(jsonOutputs, jsonData)
	}

	// Verify all JSON outputs are identical
	firstOutput := jsonOutputs[0]
	for i := 1; i < len(jsonOutputs); i++ {
		if !bytes.Equal(firstOutput, jsonOutputs[i]) {
			t.Errorf("Run %d produced different JSON output:\nFirst:\n%s\n\nRun %d:\n%s\n",
				1, prettifyJSON(firstOutput), i+1, prettifyJSON(jsonOutputs[i]))
		}
	}
}

// Test_Parser_JSON_Determinism tests that Parser marshaling is deterministic.
func Test_Parser_JSON_Determinism(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	runs := 10
	var jsonOutputs [][]byte

	for i := 0; i < runs; i++ {
		cab := fstest.MapFS{}
		p := NewParser(cab)
		p.Root = "/test/root"
		p.Section = 1

		// Add some vars in different orders
		p.Vars.Set("key1", "value1")
		p.Vars.Set("key2", "value2")
		p.Vars.Set("key3", "value3")

		jsonData, err := json.Marshal(p)
		r.NoError(err)
		jsonOutputs = append(jsonOutputs, jsonData)
	}

	// Verify all JSON outputs are identical
	firstOutput := jsonOutputs[0]
	for i := 1; i < len(jsonOutputs); i++ {
		r.Equal(firstOutput, jsonOutputs[i],
			"Parser JSON output should be deterministic")
	}
}

// Test_Element_JSON_Determinism tests that Element marshaling is deterministic.
func Test_Element_JSON_Determinism(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	runs := 10
	var jsonOutputs [][]byte

	for i := 0; i < runs; i++ {
		el := NewEl("div", nil)
		el.Set("id", "test")
		el.Set("class", "container")
		el.Set("data-foo", "bar")
		el.Set("data-baz", "qux")

		jsonData, err := json.Marshal(el)
		r.NoError(err)
		jsonOutputs = append(jsonOutputs, jsonData)
	}

	// Verify all JSON outputs are identical
	firstOutput := jsonOutputs[0]
	for i := 1; i < len(jsonOutputs); i++ {
		r.Equal(firstOutput, jsonOutputs[i],
			"Element JSON output should be deterministic")
	}
}

// Test_Snippets_JSON_Determinism tests that Snippets marshaling is deterministic.
func Test_Snippets_JSON_Determinism(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	runs := 10
	var jsonOutputs [][]byte

	for i := 0; i < runs; i++ {
		snips := &Snippets{}
		snips.init()

		// Add some rules
		snips.Add(".go", "// %s")
		snips.Add(".rb", "# %s")
		snips.Add(".js", "// %s")

		// Simulate adding snippets
		snips.mu.Lock()
		snips.snippets["file1.go"] = map[string]Snippet{
			"snippet1": {Name: "snippet1", Content: "func test() {}", Lang: "go"},
			"snippet2": {Name: "snippet2", Content: "func main() {}", Lang: "go"},
		}
		snips.snippets["file2.rb"] = map[string]Snippet{
			"snippet3": {Name: "snippet3", Content: "def test; end", Lang: "rb"},
		}
		snips.mu.Unlock()

		jsonData, err := json.Marshal(snips)
		r.NoError(err)
		jsonOutputs = append(jsonOutputs, jsonData)
	}

	// Verify all JSON outputs are identical
	firstOutput := jsonOutputs[0]
	for i := 1; i < len(jsonOutputs); i++ {
		r.Equal(firstOutput, jsonOutputs[i],
			"Snippets JSON output should be deterministic")
	}
}

// Benchmark_Document_JSON_Marshal benchmarks JSON marshaling performance
func Benchmark_Document_JSON_Marshal(b *testing.B) {
	cab := fstest.MapFS{
		"test.md": &fstest.MapFile{
			Data: []byte(`<html><body><page>
<h1 id="title" class="header">Benchmark Test</h1>
<p class="content">This is a benchmark test.</p>
<div id="container" class="wrapper">
  <h2>Section</h2>
  <p>Content here.</p>
</div>
</page></body></html>`),
		},
	}

	p := NewParser(cab)
	p.DocIDGen = func() (string, error) {
		return "bench-doc-id", nil
	}

	doc, err := p.ParseFile("test.md")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(doc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// prettifyJSON formats JSON for better error output
func prettifyJSON(data []byte) string {
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return string(data)
	}
	return out.String()
}
