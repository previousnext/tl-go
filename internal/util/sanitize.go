package util

import (
	"regexp"
	"strings"
)

var multipleSpaces = regexp.MustCompile(`\s+`)

// SanitizeComment sanitizes a comment string for display in Jira worklogs.
// It replaces newlines, carriage returns, and tabs with spaces,
// collapses multiple whitespace characters to a single space,
// and trims leading/trailing whitespace.
// Note: JSON escaping of quotes and backslashes is handled by encoding/json.
func SanitizeComment(comment string) string {
	// Replace newlines and carriage returns with space
	s := strings.ReplaceAll(comment, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")

	// Replace tabs with space
	s = strings.ReplaceAll(s, "\t", " ")

	// Collapse multiple whitespace to single space
	s = multipleSpaces.ReplaceAllString(s, " ")

	// Trim leading/trailing whitespace
	return strings.TrimSpace(s)
}
