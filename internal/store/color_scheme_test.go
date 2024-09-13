package store

import (
	"testing"
)

func TestColorSchemeValue(t *testing.T) {
	tests := []struct {
		scheme   ColorScheme
		expected int64
	}{
		{
			scheme:   Dark,
			expected: 3,
		},
		{
			scheme:   Light,
			expected: 2,
		},
		{
			scheme:   System,
			expected: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.scheme.String(), func(t *testing.T) {
			v, err := test.scheme.Value()
			if err != nil {
				t.Error(err)
			}
			if v.(int64) != test.expected {
				t.Errorf("Expected %d to equal %d", v, test.expected)
			}
		})
	}
}

func TestColorSchemeScan(t *testing.T) {
	schemes := []ColorScheme{
		System, Dark, Light,
	}

	for _, input := range schemes {
		t.Run(input.String(), func(t *testing.T) {
			v, err := input.Value()
			if err != nil {
				t.Error(err)
			}

			var scanned ColorScheme
			err = scanned.Scan(v)
			if err != nil {
				t.Error(err)
			}

			if scanned != input {
				t.Errorf("Expected %s to equal %s", scanned, input)
			}
		})
	}
}
