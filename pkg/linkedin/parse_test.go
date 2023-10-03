package linkedin

import (
	"testing"
)

func TestCleanURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://example.com/path?query=value", "https://example.com/path"},
		{"https://example.com/path", "https://example.com/path"},
		{"invalid-url", "invalid-url"},
	}

	for _, test := range tests {
		result := cleanURL(test.input)
		if result != test.expected {
			t.Errorf("For input %s, expected %s, but got %s", test.input, test.expected, result)
		}
	}
}

func TestExtractLikesCount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"1,234", 1234, false},
		{"123", 123, false},
		{"", 0, true},
	}

	for _, test := range tests {
		result, err := extractLikesCount(test.input)
		if test.hasError && err == nil {
			t.Errorf("Expected error for input %s, but got none", test.input)
		}
		if !test.hasError && result != test.expected {
			t.Errorf("For input %s, expected %d, but got %d", test.input, test.expected, result)
		}
	}
}

func TestExtractCommentsCount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"1,234 Comments", 1234, false},
		{"123 Comments", 123, false},
		{"Comments", 0, true},
	}

	for _, test := range tests {
		result, err := extractCommentsCount(test.input)
		if test.hasError && err == nil {
			t.Errorf("Expected error for input %s, but got none", test.input)
		}
		if !test.hasError && result != test.expected {
			t.Errorf("For input %s, expected %d, but got %d", test.input, test.expected, result)
		}
	}
}

func TestExtractFollowersCount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"1,234 followers", 1234, false},
		{"123 followers", 123, false},
		{"followers", 0, true},
	}

	for _, test := range tests {
		result, err := extractFollowersCount(test.input)
		if test.hasError && err == nil {
			t.Errorf("Expected error for input %s, but got none", test.input)
		}
		if !test.hasError && result != test.expected {
			t.Errorf("For input %s, expected %d, but got %d", test.input, test.expected, result)
		}
	}
}
