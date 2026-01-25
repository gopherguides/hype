package cli

import (
	"bytes"
	"flag"
	"testing"
	"time"

	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Preview_Flags(t *testing.T) {
	r := require.New(t)

	cmd := &Preview{}
	var buf bytes.Buffer

	flags, err := cmd.Flags(&buf)
	r.NoError(err)
	r.NotNil(flags)

	flags2, err := cmd.Flags(&buf)
	r.NoError(err)
	r.Equal(flags, flags2, "should return same flagset on subsequent calls")
}

func Test_Preview_Flags_Defaults(t *testing.T) {
	r := require.New(t)

	cmd := &Preview{}
	var buf bytes.Buffer

	flags, err := cmd.Flags(&buf)
	r.NoError(err)

	err = flags.Parse([]string{})
	r.NoError(err)

	r.Equal("hype.md", cmd.File)
	r.Equal(3000, cmd.Port)
	r.Equal(300*time.Millisecond, cmd.DebounceDelay)
	r.Equal("github", cmd.Theme)
	r.False(cmd.Verbose)
	r.False(cmd.OpenBrowser)
}

func Test_Preview_Flags_Custom(t *testing.T) {
	r := require.New(t)

	cmd := &Preview{}
	var buf bytes.Buffer

	flags, err := cmd.Flags(&buf)
	r.NoError(err)

	args := []string{
		"-f", "README.md",
		"-port", "8080",
		"-e", "md,go,html",
		"-w", "./src",
		"-w", "./docs",
		"-i", "**/*.md",
		"-x", "**/tmp/**",
		"-d", "500ms",
		"-v",
		"-open",
		"-theme", "github-dark",
	}

	err = flags.Parse(args)
	r.NoError(err)

	r.Equal("README.md", cmd.File)
	r.Equal(8080, cmd.Port)
	r.Equal("md,go,html", cmd.Extensions)
	r.Equal(stringSlice{"./src", "./docs"}, cmd.WatchDirs)
	r.Equal(stringSlice{"**/*.md"}, cmd.IncludeGlobs)
	r.Equal(stringSlice{"**/tmp/**"}, cmd.ExcludeGlobs)
	r.Equal(500*time.Millisecond, cmd.DebounceDelay)
	r.True(cmd.Verbose)
	r.True(cmd.OpenBrowser)
	r.Equal("github-dark", cmd.Theme)
}

func Test_Preview_Flags_Aliases(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want func(*Preview) bool
	}{
		{
			name: "watch alias",
			args: []string{"-watch", "./src"},
			want: func(p *Preview) bool { return len(p.WatchDirs) == 1 && p.WatchDirs[0] == "./src" },
		},
		{
			name: "ext alias",
			args: []string{"-ext", "md,go"},
			want: func(p *Preview) bool { return p.Extensions == "md,go" },
		},
		{
			name: "include alias",
			args: []string{"-include", "**/*.md"},
			want: func(p *Preview) bool { return len(p.IncludeGlobs) == 1 && p.IncludeGlobs[0] == "**/*.md" },
		},
		{
			name: "exclude alias",
			args: []string{"-exclude", "**/tmp/**"},
			want: func(p *Preview) bool { return len(p.ExcludeGlobs) == 1 && p.ExcludeGlobs[0] == "**/tmp/**" },
		},
		{
			name: "debounce alias",
			args: []string{"-debounce", "500ms"},
			want: func(p *Preview) bool { return p.DebounceDelay == 500*time.Millisecond },
		},
		{
			name: "verbose alias",
			args: []string{"-verbose"},
			want: func(p *Preview) bool { return p.Verbose },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)

			cmd := &Preview{}
			var buf bytes.Buffer

			flags, err := cmd.Flags(&buf)
			r.NoError(err)

			err = flags.Parse(tc.args)
			r.NoError(err)

			r.True(tc.want(cmd))
		})
	}
}

func Test_Preview_SetParser(t *testing.T) {
	r := require.New(t)

	cmd := &Preview{}
	r.Nil(cmd.Parser)

	err := cmd.SetParser(nil)
	r.NoError(err)
	r.Nil(cmd.Parser)

	var nilCmd *Preview
	err = nilCmd.SetParser(nil)
	r.Error(err)
	r.Contains(err.Error(), "preview is nil")
}

func Test_Preview_Validate(t *testing.T) {
	r := require.New(t)

	cmd := &Preview{}
	err := cmd.validate()
	r.NoError(err)

	var nilCmd *Preview
	err = nilCmd.validate()
	r.Error(err)
	r.Contains(err.Error(), "cmd is nil")
}

func Test_stringSlice_String(t *testing.T) {
	tests := []struct {
		name string
		s    stringSlice
		want string
	}{
		{
			name: "empty",
			s:    stringSlice{},
			want: "",
		},
		{
			name: "single",
			s:    stringSlice{"a"},
			want: "a",
		},
		{
			name: "multiple",
			s:    stringSlice{"a", "b", "c"},
			want: "a,b,c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			r.Equal(tc.want, tc.s.String())
		})
	}
}

func Test_stringSlice_Set(t *testing.T) {
	r := require.New(t)

	var s stringSlice

	err := s.Set("a")
	r.NoError(err)
	r.Equal(stringSlice{"a"}, s)

	err = s.Set("b")
	r.NoError(err)
	r.Equal(stringSlice{"a", "b"}, s)
}

func Test_Preview_Flags_Usage(t *testing.T) {
	r := require.New(t)

	cmd := &Preview{}
	var buf bytes.Buffer

	flags, err := cmd.Flags(&buf)
	r.NoError(err)

	err = flags.Parse([]string{"-h"})
	r.Equal(flag.ErrHelp, err)

	output := buf.String()
	r.Contains(output, "hype preview")
	r.Contains(output, "-f")
	r.Contains(output, "-port")
	r.Contains(output, "-w")
	r.Contains(output, "-e")
	r.Contains(output, "-i")
	r.Contains(output, "-x")
	r.Contains(output, "-open")
	r.Contains(output, "-theme")
}

func Test_Preview_ScopedPlugins(t *testing.T) {
	r := require.New(t)

	cmd := &Preview{}
	_ = cmd.ScopedPlugins()

	var nilCmd *Preview
	nilPlugs := nilCmd.ScopedPlugins()
	r.Nil(nilPlugs)
}

func Test_Preview_WithPlugins(t *testing.T) {
	r := require.New(t)

	cmd := &Preview{}

	err := cmd.WithPlugins(nil)
	r.Error(err)
	r.Contains(err.Error(), "fn is nil")

	err = cmd.WithPlugins(func() plugins.Plugins { return nil })
	r.NoError(err)

	var nilCmd *Preview
	err = nilCmd.WithPlugins(func() plugins.Plugins { return nil })
	r.Error(err)
	r.Contains(err.Error(), "preview is nil")
}
