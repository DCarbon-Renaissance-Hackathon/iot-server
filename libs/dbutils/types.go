package dbutils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type MapSS map[string]string

func (m *MapSS) Scan(value interface{}) error {
	if nil == m {
		m = new(MapSS)
	}
	switch vt := value.(type) {
	case string:
		return json.Unmarshal([]byte(vt), m)
	case []byte:
		return json.Unmarshal(vt, m)
	}
	return errors.New("scan value type for MapSS invalid")
}

func (m MapSS) Value() (driver.Value, error) {
	if nil == m {
		return nil, nil
	}
	return json.Marshal(m)
}

type MapSFloat map[string]float64

func (m *MapSFloat) Scan(value interface{}) error {
	if nil == m {
		m = new(MapSFloat)
	}
	switch vt := value.(type) {
	case string:
		return json.Unmarshal([]byte(vt), m)
	case []byte:
		return json.Unmarshal(vt, m)
	}
	return errors.New("scan value type for MapSFloat invalid")
}

func (m MapSFloat) Value() (driver.Value, error) {
	if nil == m {
		return nil, nil
	}
	return json.Marshal(m)
}

type Strings []string

func (ss *Strings) Scan(value interface{}) error {
	if nil == ss {
		ss = new(Strings)
	}
	switch vt := value.(type) {
	case string:
		return json.Unmarshal([]byte(vt), ss)
	case []byte:
		return json.Unmarshal(vt, ss)
	}
	return errors.New("scan value type for Strings is invalid")
}

func (ss Strings) Value() (driver.Value, error) {
	if nil == ss {
		return nil, nil
	}
	return json.Marshal(ss)
}
