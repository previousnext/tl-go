package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeComment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple comment", "Simple comment"},
		{"Line one\nLine two", "Line one Line two"},
		{"Line one\r\nLine two", "Line one Line two"},
		{"Item\tValue", "Item Value"},
		{"Too   many   spaces", "Too many spaces"},
		{"  padded  ", "padded"},
		{"A\n\n\nB", "A B"},
		{"   \n\t  ", ""},
		{"", ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, SanitizeComment(tt.input))
	}
}
