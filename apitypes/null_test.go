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
				IsSet:  true,
				String: "",
				IsNull: true,
			},
		},
		{
			input: `""`,
			expected: NullString{
				IsSet:  true,
				String: "",
				IsNull: false,
			},
		},
		{
			input: `"test"`,
			expected: NullString{
				IsSet:  true,
				String: "test",
				IsNull: false,
			},
		},
	}

	for _, table := range tables {
		var string NullString
		err := json.Unmarshal([]byte(table.input), &string)
		if err != nil {
			t.Error(err)
		}

		if string.IsSet != table.expected.IsSet {
			t.Errorf("Expected IsSet to be %t got %t", table.expected.IsSet, string.IsSet)
		}
		if string.IsNull != table.expected.IsNull {
			t.Errorf("Expected IsNull to be %t got %t", table.expected.IsNull, string.IsNull)
		}
		if string.String != table.expected.String {
			t.Errorf("Expected String to be %s got %s", table.expected.String, string.String)
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
				IsSet:  true,
				String: "",
				IsNull: true,
			},
			expect: "null",
		},
		{
			input: NullString{
				IsSet:  true,
				String: "",
				IsNull: false,
			},
			expect: `""`,
		},
		{
			input: NullString{
				IsSet:  true,
				String: "test",
				IsNull: false,
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
				String: "hello",
				IsSet:  true,
				IsNull: false,
			},
		},
		{
			name:  "Null value",
			input: nil,
			expected: NullString{
				String: "",
				IsSet:  true,
				IsNull: true,
			},
			wantErr: false,
		},
		{
			name:  "Empty string",
			input: "",
			expected: NullString{
				String: "",
				IsSet:  true,
				IsNull: false,
			},
			wantErr: false,
		},
		{
			name:  "Non-string input (integer)",
			input: 123,
			expected: NullString{
				String: "123",
				IsSet:  true,
				IsNull: false,
			},
			wantErr: false,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			var ns NullString
			err := ns.Scan(table.input)

			if table.wantErr && err == nil {
				t.Error("Expected an error but got none")
			} else if !table.wantErr && err != nil {
				t.Errorf("Expected no error, but got %s", err)
			}

			if ns.IsSet != table.expected.IsSet {
				t.Errorf("Expected IsSet to be %t got %t", table.expected.IsSet, ns.IsSet)
			}
			if ns.IsNull != table.expected.IsNull {
				t.Errorf("Expected IsNull to be %t got %t", table.expected.IsNull, ns.IsNull)
			}
			if ns.String != table.expected.String {
				t.Errorf("Expected String to be %s got %s", table.expected.String, ns.String)
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
			ns:       NullString{String: "hello", IsSet: true, IsNull: false},
			expected: "hello",
		},
		{
			name:     "Null value",
			ns:       NullString{String: "", IsSet: true, IsNull: true},
			expected: nil,
		},
		{
			name:     "Empty string",
			ns:       NullString{String: "", IsSet: true, IsNull: false},
			expected: "",
		},
		{
			name:     "Not set",
			ns:       NullString{String: "", IsSet: false, IsNull: false},
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
