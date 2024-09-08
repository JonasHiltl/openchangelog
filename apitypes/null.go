package apitypes

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Represents a nullable value.
// Supports JSON un/marshaling and implements the Scanner interface.
type NullString struct {
	str    string
	isSet  bool
	isNull bool
}

func (ns NullString) String() string {
	return ns.str
}

func (ns NullString) IsSet() bool {
	return ns.isSet
}

func (ns NullString) IsNull() bool {
	return ns.isNull
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		s.isSet = true
		s.isNull = true
		return nil
	}

	if err := json.Unmarshal(data, &s.str); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	s.isSet = true
	return nil
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if s.isNull {
		return []byte("null"), nil
	}
	return json.Marshal(s.str)
}

func (n *NullString) Scan(value interface{}) error {
	ns := sql.NullString{}
	err := ns.Scan(value)
	if err != nil {
		return err
	}

	n.str = ns.String
	n.isNull = !ns.Valid
	n.isSet = true
	return nil
}

func (n NullString) Value() (driver.Value, error) {
	ns := sql.NullString{
		String: n.str,
		Valid:  !n.isNull && n.isSet,
	}
	return ns.Value()
}
