package kstrings

import "strings"

// HasPrefix tests whether the string s begins with prefix, case-insensitive.
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

// HasSuffix tests whether the string s ends with suffix, case-insensitive.
func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix)
}
