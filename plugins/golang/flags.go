package golang

import "strings"

func CleanFlags(flags ...string) []string {
	res := make([]string, 0, len(flags))
	for _, s := range flags {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			res = append(res, s)
		}
	}
	return res
}
