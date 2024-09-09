package apitypes

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewString(t *testing.T) {
	tables := []struct {
		name    string
		input   string
		isZero  bool
		isValid bool
	}{
		{
			name:    "Empty string",
			input:   "",
			isZero:  true,
			isValid: false,
		},
		{
			name:    "Valid string",
			input:   "hello",
			isZero:  false,
			isValid: true,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			ns := NewString(table.input)
			assert.Equal(t, ns.IsZero(), table.isZero)
			assert.Equal(t, ns.IsValid(), table.isValid)
		})
	}
}

func TestUnmarshalString(t *testing.T) {
	tables := []struct {
		input    string
		expected NullString
	}{
		{
			input:    "null",
			expected: NewNullString(),
		},
		{
			input:    `""`,
			expected: NewString(""),
		},
		{
			input:    `"test"`,
			expected: NewString("test"),
		},
	}

	for _, table := range tables {
		var ns NullString
		err := json.Unmarshal([]byte(table.input), &ns)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, table.expected, ns)
	}
}

func TestMarshalString(t *testing.T) {
	tables := []struct {
		input  NullString
		expect string
	}{
		{
			input:  NewNullString(),
			expect: "null",
		},
		{
			input:  NewString(""),
			expect: `""`,
		},
		{
			input:  NewString("test"),
			expect: `"test"`,
		},
	}

	for _, table := range tables {
		output, err := json.Marshal(table.input)
		if err != nil {
			t.Error(err)
		}

		str := string(output)
		if str != table.expect {
			t.Errorf("Expected %s to equal %s", str, table.expect)
		}
	}
}

func TestScanString(t *testing.T) {
	tables := []struct {
		name     string
		input    any
		expected NullString
		wantErr  bool
	}{
		{
			name:     "Valid string",
			input:    "hello",
			expected: NewString("hello"),
		},
		{
			name:     "Null value",
			input:    nil,
			expected: NewNullString(),
			wantErr:  false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: NewString(""),
			wantErr:  false,
		},
		{
			name:     "Non-string input (integer)",
			input:    123,
			expected: NewString("123"),
			wantErr:  false,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			var ns NullString
			err := ns.Scan(table.input)

			if table.wantErr == (err == nil) {
				t.Errorf("Expected wantErr %t but got %s", table.wantErr, err)
			}

			assert.Equal(t, table.expected, ns)

		})
	}
}

func TestValueString(t *testing.T) {
	tests := []struct {
		name     string
		ns       NullString
		expected driver.Value
	}{
		{
			name:     "Valid string",
			ns:       NewString("hello"),
			expected: "hello",
		},
		{
			name:     "Null value",
			ns:       NewNullString(),
			expected: nil,
		},
		{
			name:     "Empty string",
			ns:       NewString(""),
			expected: nil,
		},
	}

	for _, table := range tests {
		t.Run(table.name, func(t *testing.T) {
			value, err := table.ns.Value()
			if err != nil {
				t.Errorf("Did not expect error %s", err)
			}
			if value != table.expected {
				t.Errorf("Expected %s to equal %s", value, table.expected)
			}
		})
	}
}
