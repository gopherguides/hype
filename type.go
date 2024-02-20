package hype

import (
	"fmt"
	"strings"
)

func toType(x any) string {
	s := fmt.Sprintf("%T", x)
	return strings.TrimPrefix(s, "*")
}
