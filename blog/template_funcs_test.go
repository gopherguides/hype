package blog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceItems(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}

	t.Run("first 3", func(t *testing.T) {
		result, err := sliceItems(items, 0, 3)
		require.NoError(t, err)
		require.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("first more than length", func(t *testing.T) {
		result, err := sliceItems(items, 0, 10)
		require.NoError(t, err)
		require.Equal(t, items, result)
	})

	t.Run("after 2", func(t *testing.T) {
		result, err := sliceItems(items, 2, -1)
		require.NoError(t, err)
		require.Equal(t, []string{"c", "d", "e"}, result)
	})

	t.Run("after more than length", func(t *testing.T) {
		result, err := sliceItems(items, 10, -1)
		require.NoError(t, err)
		require.Equal(t, []string{}, result)
	})

	t.Run("first 0", func(t *testing.T) {
		result, err := sliceItems(items, 0, 0)
		require.NoError(t, err)
		require.Equal(t, []string{}, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result, err := sliceItems([]string{}, 0, 3)
		require.NoError(t, err)
		require.Equal(t, []string{}, result)
	})

	t.Run("works with Article slice", func(t *testing.T) {
		articles := []Article{
			{Title: "First"},
			{Title: "Second"},
			{Title: "Third"},
		}
		result, err := sliceItems(articles, 0, 2)
		require.NoError(t, err)
		got := result.([]Article)
		require.Len(t, got, 2)
		require.Equal(t, "First", got[0].Title)
		require.Equal(t, "Second", got[1].Title)
	})

	t.Run("rejects non-slice", func(t *testing.T) {
		_, err := sliceItems("not a slice", 0, 3)
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected a slice")
	})
}

func TestListPagesPathValidation(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid simple path", "blog", false},
		{"valid nested path", "news/archive", false},
		{"absolute path", "/tmp/evil", true},
		{"traversal path", "../escape", true},
		{"dot-dot in middle", "foo/../../etc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateListPagePath(tt.path)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid listPages path")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
