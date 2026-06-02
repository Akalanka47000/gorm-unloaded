package main

import (
	"database/sql/driver"
	"fmt"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusActive    Status = "active"
	StatusCancelled Status = "cancelled"
)

// Value implements driver.Valuer — called by database/sql when writing to the DB.
// Returns the underlying string so PostgreSQL stores it as text.
func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements sql.Scanner — called by database/sql when reading from the DB.
// Accepts the raw database value and populates the receiver.
func (s *Status) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into Status", value)
	}
	*s = Status(str)
	return nil
}
