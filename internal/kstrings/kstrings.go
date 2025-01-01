package kstrings

import (
	"strings"
)

// HasPrefix tests whether the string s begins with prefix, case-insensitive.
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

// HasSuffix tests whether the string s ends with suffix, case-insensitive.
func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix)
}

// ToSnakeCase converts a string into snake_case.
func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, char := range str {
		if char >= 'A' && char <= 'Z' {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(char + 'a' - 'A')
		} else if char == ' ' {
			result.WriteRune('_')
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// ToCamelCase converts a string into camelCase.
func ToCamelCase(str string) string {
	var result strings.Builder
	capitalizeNext := false

	for i, char := range str {
		if char == ' ' || char == '_' {
			capitalizeNext = true
		} else {
			if capitalizeNext {
				result.WriteRune(char - 'a' + 'A')
				capitalizeNext = false
			} else {
				if i == 0 && char >= 'A' && char <= 'Z' {
					result.WriteRune(char + 'a' - 'A')
				} else {
					result.WriteRune(char)
				}
			}
		}
	}

	return result.String()
}
