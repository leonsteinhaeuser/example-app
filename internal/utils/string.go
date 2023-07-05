package utils

import "strings"

func ReplaceAllToReplace(s string, new string, toReplace ...string) string {
	for _, r := range toReplace {
		s = strings.ReplaceAll(s, r, new)
	}
	return s
}
