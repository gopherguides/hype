package hype

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DocType_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	dt := &DocType{
		Node: NewNode(DocTypeNode(t, "html8")),
	}

	exp := "<!doctype html8>\n"
	act := dt.String()
	r.Equal(exp, act)
}

func Test_DocType_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	dt := &DocType{
		Node: NewNode(DocTypeNode(t, "html8")),
	}

	exp := `{"data":"html8","type":"doctype"}`
	b, err := json.Marshal(dt)
	r.NoError(err)
	act := string(b)
	r.Equal(exp, act)
}
