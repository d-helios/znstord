package zfs

import "fmt"

// Error - common error structure
type Error struct {
	Err    error
	Debug  string
	Stderr string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %q => %s", e.Err, e.Debug, e.Stderr)
}
