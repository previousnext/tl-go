package util

import "strings"

// AbbreviateProjectCategory takes a project category and returns first character from each word.
func AbbreviateProjectCategory(category string) string {
	if category == "" {
		return "N/A"
	}
	// splits a string into words based on spaces.
	words := strings.Fields(category)
	var initials string
	for _, word := range words {
		if len(word) > 0 {
			initials += string(word[0])
		}
	}
	return initials
}
