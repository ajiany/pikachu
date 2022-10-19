package httpx

import "strings"

func Escape(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\\{", "{")
	s = strings.ReplaceAll(s, "\\[", "[")
	return s
}

func Truncation(str string, length int) string {
	if len(str) <= length {
		return str
	}

	return str[0:length]
}

func DefaultString(str, other string) string {
	if str == "" {
		return other
	}

	return str
}
