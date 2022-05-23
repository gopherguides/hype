package module

import (
	"os"
	"testing"

	"github.com/gopherguides/hypewriter"
	"github.com/gopherguides/hypewriter/hyper"
	"github.com/stretchr/testify/require"
)

func Test_Module(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	cab := os.DirFS(".")

	p, err := hyper.NewParser(cab, ".")
	r.NoError(err)

	bind, err := hypewriter.NewBinder("Test", p)
	r.NoError(err)

	mods, err := bind.BindAll()
	r.NoError(err)
	r.Len(mods, 1)

}
