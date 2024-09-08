package apitypes

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

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

		if ns.IsValid() != table.expected.IsValid() {
			t.Errorf("Expected isValid to be %t got %t", table.expected.IsValid(), ns.IsValid())
		}
		if ns.IsNull() != table.expected.IsNull() {
			t.Errorf("Expected isNull to be %t got %t", table.expected.IsNull(), ns.IsNull())
		}
		if ns.String() != table.expected.String() {
			t.Errorf("Expected str to be %s got %s", table.expected.String(), ns.String())
		}
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

			if ns.IsValid() != table.expected.IsValid() {
				t.Errorf("Expected IsValid to be %t got %t", table.expected.IsValid(), ns.IsValid())
			}
			if ns.IsNull() != table.expected.IsNull() {
				t.Errorf("Expected isNull to be %t got %t", table.expected.IsNull(), ns.IsNull())
			}
			if ns.String() != table.expected.String() {
				t.Errorf("Expected str to be %s got %s", table.expected.String(), ns.String())
			}
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
			expected: "",
		},
		{
			name:     "Not set",
			ns:       NullString{str: nil, isNull: false},
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
