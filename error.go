package pushover

import (
	"fmt"
	"net"
)

// TemporaryError implementes net.Error and represents temporary error.
// Request may be retried later after 5 second delay.
type TemporaryError struct {
	StatusCode int
	Message    string
}

// FatalError implementes net.Error and represents fatal error.
// Request should not be retried.
type FatalError struct {
	StatusCode int
	Message    string
}

// check interfaces
var (
	_ net.Error = &TemporaryError{}
	_ net.Error = &FatalError{}
)

func (c *TemporaryError) Error() string {
	return fmt.Sprintf("pushover: temporary error: %s (%d)", c.Message, c.StatusCode)
}
func (c *TemporaryError) Temporary() bool { return true }
func (c *TemporaryError) Timeout() bool   { return false }

func (c *FatalError) Error() string {
	return fmt.Sprintf("pushover: fatal error: %s (%d)", c.Message, c.StatusCode)
}
func (c *FatalError) Temporary() bool { return false }
func (c *FatalError) Timeout() bool   { return false }
