package utils

import (
	"fmt"
	"strings"
)

// FormatType is a custom type for the format.
type FormatType string

const (
	JSON FormatType = "json"
	CSV  FormatType = "csv"
)

// String returns the string representation of the FormatType.
func (f FormatType) String() string {
	return string(f)
}

// SetFormat sets the format from a string, handling case-insensitivity.
func SetFormat(s string) (FormatType, error) {
	switch strings.ToLower(s) {
	case "json":
		return JSON, nil
	case "csv":
		return CSV, nil
	default:
		return "", fmt.Errorf("unknown format: %s", s)
	}
}
