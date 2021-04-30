package pushover

import "net"

type Error struct {
	Err error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Timeout() bool {
	if ne, ok := e.Err.(net.Error); ok {
		return ne.Timeout()
	}

	return false
}

func (e *Error) Temporary() bool {
	if ne, ok := e.Err.(net.Error); ok {
		return ne.Temporary()
	}

	return false
}

// check interfaces
var (
	_ error     = (*Error)(nil)
	_ net.Error = (*Error)(nil)
)
