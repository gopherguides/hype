package hype

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Body_JSON(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	bd := AttrNode(t, "body", Attributes{
		"id": "main",
	})
	bd.FirstChild = TextNode(t, "hi")

	body := &Body{
		Node: NewNode(bd),
	}

	b, err := json.Marshal(body)
	r.NoError(err)

	exp := `{"atom":"body","attributes":{"id":"main"},"children":[{"data":"hi","type":"text"}],"data":"body","type":"element"}`

	act := string(b)

	r.Equal(exp, act)
}
