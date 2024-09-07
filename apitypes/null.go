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
	String string
	IsSet  bool
	IsNull bool
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		s.IsSet = true
		s.IsNull = true
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	s.IsSet = true
	return nil
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if s.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (n *NullString) Scan(value interface{}) error {
	ns := sql.NullString{}
	err := ns.Scan(value)
	if err != nil {
		return err
	}

	n.String = ns.String
	n.IsNull = !ns.Valid
	n.IsSet = true
	return nil
}

func (n NullString) Value() (driver.Value, error) {
	ns := sql.NullString{
		String: n.String,
		Valid:  !n.IsNull && n.IsSet,
	}
	return ns.Value()
}
