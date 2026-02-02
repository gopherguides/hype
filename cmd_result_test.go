package hype

import (
	"io"
	"testing"
	"time"

	"github.com/markbates/clam"
)

func Test_applyReplacements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		attrs    map[string]string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name: "replace go version",
			attrs: map[string]string{
				"replace-1":      `go1\.\d+\.\d+`,
				"replace-1-with": "goX.X.X",
			},
			input:    "go version go1.21.5 darwin/arm64",
			expected: "go version goX.X.X darwin/arm64",
		},
		{
			name: "replace timestamp",
			attrs: map[string]string{
				"replace-1":      `\d{4}-\d{2}-\d{2}`,
				"replace-1-with": "[DATE]",
			},
			input:    "Created on 2024-01-15",
			expected: "Created on [DATE]",
		},
		{
			name: "multiple replacements in order",
			attrs: map[string]string{
				"replace-1":      `go1\.\d+\.\d+`,
				"replace-1-with": "goX.X.X",
				"replace-2":      `\d{4}-\d{2}-\d{2}`,
				"replace-2-with": "[DATE]",
			},
			input:    "go1.21.5 built on 2024-01-15",
			expected: "goX.X.X built on [DATE]",
		},
		{
			name: "replacement to empty string",
			attrs: map[string]string{
				"replace-1":      ` \(.*?\)`,
				"replace-1-with": "",
			},
			input:    "Hello (remove this) World",
			expected: "Hello World",
		},
		{
			name: "no matching pattern",
			attrs: map[string]string{
				"replace-1":      `[0-9a-f]{8}-.*`,
				"replace-1-with": "[UUID]",
			},
			input:    "No UUID here",
			expected: "No UUID here",
		},
		{
			name: "invalid regex returns error",
			attrs: map[string]string{
				"replace-1":      `[invalid`,
				"replace-1-with": "x",
			},
			input:   "some text",
			wantErr: true,
		},
		{
			name: "missing with attribute defaults to empty",
			attrs: map[string]string{
				"replace-1": `remove-me`,
			},
			input:    "please remove-me now",
			expected: "please  now",
		},
		{
			name:     "no replacements",
			attrs:    map[string]string{},
			input:    "unchanged text",
			expected: "unchanged text",
		},
		{
			name: "uuid replacement",
			attrs: map[string]string{
				"replace-1":      `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
				"replace-1-with": "[UUID]",
			},
			input:    "ID: 550e8400-e29b-41d4-a716-446655440000",
			expected: "ID: [UUID]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ats := &Attributes{}
			for k, v := range tc.attrs {
				ats.Set(k, v)
			}
			result, err := applyReplacements(ats, tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if result != tc.expected {
				t.Errorf("got %q, want %q", result, tc.expected)
			}
		})
	}
}

func Test_CmdResult_MarshalJSON(t *testing.T) {
	t.Parallel()

	cr := &CmdResult{
		Element: NewEl("cmd", nil),
		Result: &clam.Result{
			Args:     []string{"echo", "hello"},
			Dir:      "/tmp",
			Duration: time.Second,
			Env:      []string{"FOO=bar", "BAR=baz"},
			Err:      io.EOF,
			Exit:     1,
			Stderr:   []byte("nothing"),
			Stdout:   []byte("foo\nbar\nbaz\n"),
		},
	}

	testJSON(t, "cmd_result", cr)

}
