package blog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsBuiltinTheme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"suspended", "suspended", true},
		{"developer", "developer", true},
		{"cards", "cards", true},
		{"unknown", "unknown", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tt.want, IsBuiltinTheme(tt.in))
		})
	}
}

func TestListBuiltinThemes(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	themes := ListBuiltinThemes()
	r.Len(themes, 3)
	r.Contains(themes, "suspended")
	r.Contains(themes, "developer")
	r.Contains(themes, "cards")

	themes[0] = "modified"
	r.NotEqual("modified", ListBuiltinThemes()[0])
}

func TestParseYAML(t *testing.T) {
	t.Parallel()

	t.Run("valid theme YAML", func(t *testing.T) {
		r := require.New(t)
		data := []byte("name: My Theme\ndescription: A theme\nauthor: Test\nversion: 1.0\nrepository: https://example.com\npreview: https://preview.com\n")
		info := &ThemeInfo{}
		err := parseYAML(data, info)
		r.NoError(err)
		r.Equal("My Theme", info.Name)
		r.Equal("A theme", info.Description)
		r.Equal("Test", info.Author)
		r.Equal("1.0", info.Version)
		r.Equal("https://example.com", info.Repository)
		r.Equal("https://preview.com", info.Preview)
	})

	t.Run("skips comments and blank lines", func(t *testing.T) {
		r := require.New(t)
		data := []byte("# comment\n\nname: Theme\n")
		info := &ThemeInfo{}
		err := parseYAML(data, info)
		r.NoError(err)
		r.Equal("Theme", info.Name)
	})

	t.Run("handles quoted values", func(t *testing.T) {
		r := require.New(t)
		data := []byte("name: \"Quoted Theme\"\nauthor: 'Single Quoted'\n")
		info := &ThemeInfo{}
		err := parseYAML(data, info)
		r.NoError(err)
		r.Equal("Quoted Theme", info.Name)
		r.Equal("Single Quoted", info.Author)
	})
}
