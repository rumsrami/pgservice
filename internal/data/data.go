package data

import (
	"github.com/pkg/errors"
	"time"
)

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("not found")
)

// TimeNowUTC returns current utc time
func TimeNowUTC() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
