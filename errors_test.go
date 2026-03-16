package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrIsNil_Error(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	e := ErrIsNil("parser")
	r.Equal("parser is nil", e.Error())
}

func TestWrapNodeErr(t *testing.T) {
	t.Parallel()

	t.Run("nil error returns nil", func(t *testing.T) {
		r := require.New(t)
		r.NoError(WrapNodeErr(Text("x"), nil))
	})

	t.Run("wraps non-Tag node with type", func(t *testing.T) {
		r := require.New(t)
		err := WrapNodeErr(Text("x"), errors.New("boom"))
		r.Error(err)
		r.Contains(err.Error(), "hype.Text")
		r.Contains(err.Error(), "boom")
	})
}

type jsonErr struct {
	Msg string `json:"msg"`
}

func (e jsonErr) Error() string { return e.Msg }
func (e jsonErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"msg": e.Msg})
}

func TestErrForJSON(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		r := require.New(t)
		r.Nil(errForJSON(nil))
	})

	t.Run("json.Marshaler returned as-is", func(t *testing.T) {
		r := require.New(t)
		je := jsonErr{Msg: "test"}
		result := errForJSON(je)
		r.Equal(je, result)
	})

	t.Run("regular error returns string", func(t *testing.T) {
		r := require.New(t)
		result := errForJSON(fmt.Errorf("simple"))
		r.Equal("simple", result)
	})
}

func TestToError(t *testing.T) {
	t.Parallel()

	t.Run("nil returns empty string", func(t *testing.T) {
		r := require.New(t)
		r.Empty(toError(nil))
	})

	t.Run("json.Marshaler returns JSON", func(t *testing.T) {
		r := require.New(t)
		result := toError(jsonErr{Msg: "test"})
		r.Contains(result, "test")
	})

	t.Run("regular error returns Error()", func(t *testing.T) {
		r := require.New(t)
		r.Equal("boom", toError(fmt.Errorf("boom")))
	})
}
