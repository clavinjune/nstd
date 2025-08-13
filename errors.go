package nstd

import (
	"context"
	"fmt"
	"os"
)

var (
	_ error        = (*shutdownCause)(nil)
	_ fmt.Stringer = (*shutdownCause)(nil)
	_ os.Signal    = (*shutdownCause)(nil)
)

// shutdownCause is used in graceful shutdown to shows which signal triggers the graceful shutdown
type shutdownCause struct {
	s os.Signal
}

// newShutdownCause creates a new shutdownCause with the given signal.
func newShutdownCause(s os.Signal) error {
	return &shutdownCause{s: s}
}

// Error returns the signal that causes
func (s *shutdownCause) Error() string {
	return fmt.Sprintf("%s: %s", context.Canceled, s.s.String())
}

// Unwrap returns the underlying error, which is context.Canceled.
func (s *shutdownCause) Unwrap() error {
	return context.Canceled
}

// Signal implements the os.Signal interface, allowing shutdownCause to be used as a signal.
func (s *shutdownCause) Signal() {
}

// String returns a string representation of the shutdownCause.
func (s *shutdownCause) String() string {
	return s.s.String()
}
