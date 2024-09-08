package apitypes

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Represents a value that can be null, not set, or set
// Supports JSON un/marshaling and implements the Scanner interface.
type NullString struct {
	str string
	// whether the value is null. JSON null value.
	isNull bool
}

// Create a new NullString with a valid value
func NewString(str string) NullString {
	return NullString{str: str}
}

// Creates a new null NullString
func NewNullString() NullString {
	return NullString{isNull: true}
}

// Returns "" if NullString is null or not valid, else the value.
func (ns NullString) String() string {
	if ns.IsNull() {
		return ""
	}

	return ns.str
}

// Returns true if the string is defined, otherwise false.
func (ns NullString) IsZero() bool {
	return ns.str == ""
}

// Returns true if the string is null, otherwiese false.
func (ns NullString) IsNull() bool {
	return ns.isNull
}

// Returns true if ns is neither null or zero value.
func (ns NullString) IsValid() bool {
	return !ns.IsNull() && !ns.IsZero()
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		ns.isNull = true
		return nil
	}

	if err := json.Unmarshal(data, &ns.str); err != nil {
		return fmt.Errorf("NullString: couldn't unmarshal JSON: %w", err)
	}

	return nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.IsNull() {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String())
}

func (n *NullString) Scan(value interface{}) error {
	ns := sql.NullString{}
	err := ns.Scan(value)
	if err != nil {
		return err
	}

	n.isNull = !ns.Valid // !valid means value is NULL in db
	if ns.Valid {
		n.str = ns.String
	}
	return nil
}

func (n NullString) Value() (driver.Value, error) {
	ns := sql.NullString{
		String: n.String(),
		Valid:  n.IsValid(),
	}
	return ns.Value()
}
