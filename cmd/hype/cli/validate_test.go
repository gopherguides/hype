package cli

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Validate_Main_Clean(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/validate/valid")
	r.NoError(err)

	cmd := &Validate{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "module.md"})
	r.NoError(err)
}

func Test_Validate_Main_WithErrors(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/validate/errors")
	r.NoError(err)

	cmd := &Validate{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "module.md"})
	r.Error(err)
	r.Contains(err.Error(), "validation failed")
}

func Test_Validate_Main_JSONFormat(t *testing.T) {
	r := require.New(t)

	pwd, err := filepath.Abs("testdata/validate/errors")
	r.NoError(err)

	cmd := &Validate{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = cmd.Main(ctx, pwd, []string{"-f", "module.md", "-format", "json"})
	r.Error(err)
	r.Contains(err.Error(), "validation failed")
}
