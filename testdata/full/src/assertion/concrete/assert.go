package assert

import (
	"bytes"
	"io"
	"time"
)

// snippet: def
func WriteNow(w io.Writer) error {
	now := time.Now()

	if bb, ok := w.(*bytes.Buffer); ok {
		bb.WriteString(now.String())
		return nil
	}

	w.Write([]byte(now.String()))

	return nil
}

// snippet: def
