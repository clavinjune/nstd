package nstd

import (
	"errors"
	"testing"
)

// RequireErrIs checks if the error is of a specific type using errors.Is, and if not, it fails the test with a formatted message.
func RequireErrIs(tb testing.TB, got error, want error) {
	tb.Helper()
	if !errors.Is(got, want) {
		tb.Errorf("got: %+v, want: %+v", got, want)
	}
}

// RequireErrAs checks if the error can be cast to a specific type using errors.As, and if not, it fails the test with a formatted message.
func RequireErrAs(tb testing.TB, got error, want any) {
	tb.Helper()
	if !errors.As(got, want) {
		tb.Errorf("got: %+v, want: %+v", got, want)
	}
}

// RequireEqual checks if the expected and actual values are equal, and if not, it fails the test with a formatted message.
func RequireEqual[T comparable](tb testing.TB, got, want T) {
	tb.Helper()

	if got != want {
		tb.Errorf("got: %+v, want: %+v", got, want)
	}
}

// RequireTrue checks if the condition is true, and if not, it fails the test with a message.
func RequireTrue(tb testing.TB, got bool) {
	tb.Helper()

	if !got {
		tb.Errorf("got: false, want: true")
	}
}

// RequireNotNil checks if the object is not nil, and if it is nil, it fails the test with a message.
func RequireNotNil(tb testing.TB, o any) {
	tb.Helper()

	if o == nil {
		tb.Error("got: nil, want: not nil")
	}
}

// RequireNil checks if the object is nil, and if it is not nil, it fails the test with a message.
func RequireNil(tb testing.TB, o any) {
	tb.Helper()

	if o != nil {
		tb.Errorf("got: %+v, want: nil", o)
	}
}
