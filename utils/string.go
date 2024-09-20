package utils

import "strings"

func HasStringKey(text, key string) bool {
	return strings.Contains(text, key)
}
