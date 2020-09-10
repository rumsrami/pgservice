package data

import (
	"time"
)

// Info represents a title
type Info struct {
	Title     string    `db:"title" json:"Title"`
	UUID4     string    `db:"uuid4" json:"UUID4"`
	Timestamp time.Time `db:"timestamp" json:"Timestamp"`
}
