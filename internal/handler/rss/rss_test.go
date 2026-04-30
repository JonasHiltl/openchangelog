package rss

import (
	"testing"
	"time"
)

func TestToPubDate(t *testing.T) {
	// Define a fixed location for consistency if needed, 
	// though toPubDate handles UTC to GMT conversion.
	utc := time.UTC

	testCases := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "Standard UTC date",
			input:    time.Date(2026, time.March, 20, 16, 24, 0, 0, utc),
			expected: "Fri, 20 Mar 2026 16:24:00 GMT",
		},
		{
			name:     "Leap year date",
			input:    time.Date(2024, time.February, 29, 0, 0, 0, 0, utc),
			expected: "Thu, 29 Feb 2024 00:00:00 GMT",
		},
		{
			name:     "Zero value date",
			input:    time.Time{},
			expected: "Mon, 01 Jan 0001 00:00:00 GMT",
		},
		{
			name:     "Ensure UTC is replaced by GMT",
			// time.RFC1123 outputs "UTC" for UTC location
			input:    time.Date(2023, time.October, 27, 10, 0, 0, 0, utc),
			expected: "Fri, 27 Oct 2023 10:00:00 GMT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toPubDate(tc.input)
			if result != tc.expected {
				t.Errorf("toPubDate() = %q, expected %q", result, tc.expected)
			}
		})
	}
}
