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

	switch i.(type) {
	case *bytes.Buffer:
		fmt.Println("type was a *bytes.Buffer", now)
	case io.StringWriter:
		fmt.Println("type was a io.StringWriter", now)
	case io.Writer:
		fmt.Println("type was a io.Writer", now)
	}

	return fmt.Errorf("can not write to %T", i)
}

// snippet: def
