package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var (
	errMapInvalidSource = errors.New("map requires []byte slice")
	errMapInvalid       = errors.New("map type assertion failed")
)

// Map implements a database-compatible map
type Map map[string]interface{}

// NewMap returns a new map
func NewMap() Map {
	return Map{}
}

func (m Map) GetInt(key string) int {
	return int(m.GetFloat64(key))
}

func (m Map) GetFloat64(key string) float64 {
	val, _ := m[key].(float64)
	return val
}

func (m Map) GetString(key string) string {
	val, _ := m[key].(string)
	return val
}

func (m Map) GetBool(key string) bool {
	val, _ := m[key].(bool)
	return val
}

// Value returns the db driver value
func (m Map) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan scants the given value into the map
func (m *Map) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errMapInvalidSource
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*m, ok = i.(map[string]interface{})
	if !ok {
		return errMapInvalid
	}

	return nil
}
