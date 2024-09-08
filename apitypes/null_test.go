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
			input: "null",
			expected: NullString{
				isSet:  true,
				str:    "",
				isNull: true,
			},
		},
		{
			input: `""`,
			expected: NullString{
				isSet:  true,
				str:    "",
				isNull: false,
			},
		},
		{
			input: `"test"`,
			expected: NullString{
				isSet:  true,
				str:    "test",
				isNull: false,
			},
		},
	}

	for _, table := range tables {
		var string NullString
		err := json.Unmarshal([]byte(table.input), &string)
		if err != nil {
			t.Error(err)
		}

		if string.isSet != table.expected.isSet {
			t.Errorf("Expected isSet to be %t got %t", table.expected.isSet, string.isSet)
		}
		if string.isNull != table.expected.isNull {
			t.Errorf("Expected isNull to be %t got %t", table.expected.isNull, string.isNull)
		}
		if string.str != table.expected.str {
			t.Errorf("Expected str to be %s got %s", table.expected.str, string.str)
		}
	}
}

func TestMarshalString(t *testing.T) {
	tables := []struct {
		input  NullString
		expect string
	}{
		{
			input: NullString{
				isSet:  true,
				str:    "",
				isNull: true,
			},
			expect: "null",
		},
		{
			input: NullString{
				isSet:  true,
				str:    "",
				isNull: false,
			},
			expect: `""`,
		},
		{
			input: NullString{
				isSet:  true,
				str:    "test",
				isNull: false,
			},
			expect: `"test"`,
		},
	}

	for _, table := range tables {
		output, err := json.Marshal(table.input)
		if err != nil {
			t.Error(err)
		}

		string := string(output)
		if string != table.expect {
			t.Errorf("Expected %s to equal %s", string, table.expect)
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
			name:  "Valid string",
			input: "hello",
			expected: NullString{
				str:    "hello",
				isSet:  true,
				isNull: false,
			},
		},
		{
			name:  "Null value",
			input: nil,
			expected: NullString{
				str:    "",
				isSet:  true,
				isNull: true,
			},
			wantErr: false,
		},
		{
			name:  "Empty string",
			input: "",
			expected: NullString{
				str:    "",
				isSet:  true,
				isNull: false,
			},
			wantErr: false,
		},
		{
			name:  "Non-string input (integer)",
			input: 123,
			expected: NullString{
				str:    "123",
				isSet:  true,
				isNull: false,
			},
			wantErr: false,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			var ns NullString
			err := ns.Scan(table.input)

			if table.wantErr == (err == nil) {
				t.Errorf("Expected wantErr %t but got %s", table.wantErr, err)
			}

			if ns.isSet != table.expected.isSet {
				t.Errorf("Expected isSet to be %t got %t", table.expected.isSet, ns.isSet)
			}
			if ns.isNull != table.expected.isNull {
				t.Errorf("Expected isNull to be %t got %t", table.expected.isNull, ns.isNull)
			}
			if ns.str != table.expected.str {
				t.Errorf("Expected str to be %s got %s", table.expected.str, ns.str)
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
			ns:       NullString{str: "hello", isSet: true, isNull: false},
			expected: "hello",
		},
		{
			name:     "Null value",
			ns:       NullString{str: "", isSet: true, isNull: true},
			expected: nil,
		},
		{
			name:     "Empty string",
			ns:       NullString{str: "", isSet: true, isNull: false},
			expected: "",
		},
		{
			name:     "Not set",
			ns:       NullString{str: "", isSet: false, isNull: false},
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
