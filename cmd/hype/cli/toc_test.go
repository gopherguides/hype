package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_TOC_Main(t *testing.T) {
	t.Skip()
	t.Parallel()

	r := require.New(t)

	root := "testdata/toc"

	cmd := &TOC{}

	out := &bytes.Buffer{}
	cmd.Out = out

	args := []string{root}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := cmd.Main(ctx, root, args)
	r.NoError(err)

	exp := `TODO`
	act := out.String()
	act = strings.TrimSpace(act)

	r.Equal(exp, act)
}
