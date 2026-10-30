// Package kerr defines the kernel's error shape. Every command that cannot run (bad usage, missing input,
// I/O failure, broken precondition) exits 2 and prints {"error":{"kind":…,"message":…}} on stderr.
// A verdict (an invalid object, a failed gate, a refused merge) is not an error: the command prints its
// report on stdout and exits 1.
package kerr

import "fmt"

// Exit codes shared by every command.
const (
	ExitOK      = 0
	ExitVerdict = 1 // invalid / gate failed / refused: report on stdout
	ExitError   = 2 // usage, I/O, precondition: error JSON on stderr
)

// Error is a kernel error with a stable kind.
type Error struct {
	Kind    string         `json:"kind"`
	Message string         `json:"message"`
	Detail  map[string]any `json:"detail,omitempty"`
}

func (e *Error) Error() string { return e.Kind + ": " + e.Message }

// New makes an Error.
func New(kind, format string, a ...any) *Error {
	return &Error{Kind: kind, Message: fmt.Sprintf(format, a...)}
}

// With returns a copy of e carrying detail.
func (e *Error) With(detail map[string]any) *Error {
	c := *e
	c.Detail = detail
	return &c
}

// Wrap turns any error into an *Error of the given kind (keeping an existing *Error unchanged). A nil err
// gives a nil error interface (never a typed nil pointer).
func Wrap(kind string, err error) error {
	if err == nil {
		return nil
	}
	if ke, ok := err.(*Error); ok {
		return ke
	}
	return &Error{Kind: kind, Message: err.Error()}
}
