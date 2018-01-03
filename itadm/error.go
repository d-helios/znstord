package itadm

import "fmt"

// Error structure
type Error struct {
	Err    error
	Debug  string
	Stderr string
}

// returns the string representation of an Error.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %q => %s", e.Err, e.Debug, e.Stderr)
}
