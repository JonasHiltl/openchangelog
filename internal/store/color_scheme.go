package store

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/jonashiltl/openchangelog/apitypes"
)

type ColorScheme int

const (
	Automatic ColorScheme = 1
	Light     ColorScheme = 2
	Dark      ColorScheme = 3
)

func NewColorScheme(cs apitypes.ColorScheme) ColorScheme {
	switch cs {
	case apitypes.Automatic:
		return Automatic
	case apitypes.Dark:
		return Dark
	case apitypes.Light:
		return Light
	}
	return 0
}

func (cs ColorScheme) String() string {
	switch cs {
	case Automatic:
		return "automatic"
	case Light:
		return "light"
	case Dark:
		return "dark"
	}
	return "unkown"
}

func (cs *ColorScheme) Scan(value interface{}) error {
	i, ok := value.(int64)
	if !ok {
		return errors.New("ColorScheme.Scan: value is not an int64")
	}

	switch ColorScheme(i) {
	case Automatic, Light, Dark:
		*cs = ColorScheme(i)
		return nil
	default:
		return fmt.Errorf("ColorScheme.Scan: failed to scan %d", i)
	}
}

func (cs ColorScheme) Value() (driver.Value, error) {
	return int64(cs), nil
}
