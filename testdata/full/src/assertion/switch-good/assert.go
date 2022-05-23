package assert

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// snippet: def
func WriteNow(i any) error {

	now := time.Now().String()

	switch t := i.(type) {
	case *bytes.Buffer:
		t.WriteString(now)
	case io.StringWriter:
		t.WriteString(now)
	case io.Writer:
		t.Write([]byte(now))
	}

	return fmt.Errorf("can not write to %T", i)
}

// snippet: def
