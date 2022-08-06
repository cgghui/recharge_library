package sys

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
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
	var s string
	err := json.Unmarshal(b, &s)
	if err == nil {
		n.Time, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
		if err == nil {
			n.Valid = true
		}
	}
	return err
}
