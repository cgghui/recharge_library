package sys

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Time sql.NullTime

func (n *Time) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

func (n Time) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n Time) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time.Format("2006-01-02 15:04:05"))
	}
	return json.Marshal(nil)
}

func (n *Time) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Time)
	if err == nil {
		n.Time = n.Time.Local()
		n.Valid = true
	}
	return err
}
