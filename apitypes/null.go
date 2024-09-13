package apitypes

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Represents a nullable value.
// Supports JSON and SQL un/marshaling.
type Null[T comparable] struct {
	v      T
	isNull bool
}

func NewValue[T comparable](v T) Null[T] {
	return Null[T]{
		v: v,
	}
}

func NewNull[T comparable]() Null[T] {
	return Null[T]{
		isNull: true,
	}
}

func (n Null[T]) IsZero() bool {
	if n.IsNull() {
		return false
	}
	return n.v == *new(T)
}

func (n Null[T]) IsNull() bool {
	return n.isNull
}

// Returns true if n is neither null or zero value.
func (n Null[T]) IsValid() bool {
	return !n.IsNull() && !n.IsZero()
}

// Returns zero value if n is null, else it's internal value.
func (n Null[T]) V() T {
	if n.IsNull() {
		return *new(T)
	}

	return n.v
}

func (n *Null[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		n.isNull = true
		return nil
	}

	if err := json.Unmarshal(data, &n.v); err != nil {
		return fmt.Errorf("Null: couldn't unmarshal JSON: %w", err)
	}

	return nil
}

func (ns Null[T]) MarshalJSON() ([]byte, error) {
	if ns.IsNull() {
		return []byte("null"), nil
	}
	return json.Marshal(ns.V())
}

func (n *Null[T]) Scan(value interface{}) error {
	sn := sql.Null[T]{}
	err := sn.Scan(value)
	if err != nil {
		return err
	}

	n.isNull = !sn.Valid // !valid means value is NULL in db
	if sn.Valid {
		n.v = sn.V
	}
	return nil
}

func (n Null[T]) Value() (driver.Value, error) {
	sn := sql.Null[T]{
		V:     n.V(),
		Valid: n.IsValid(), // this way we also store zero values as NULL in db
	}
	return sn.Value()
}

type NullString = Null[string]

// Create a new NullString from a string value
func NewString(str string) NullString {
	return NewValue(str)
}

// Creates a new null NullString
func NewNullString() NullString {
	return NewNull[string]()
}

type NullColorScheme = Null[ColorScheme]
