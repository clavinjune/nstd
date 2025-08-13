package nstd

import (
	"strings"
	"testing"
)

// RequireNoErr checks if the error is nil, and if not, it fails the test with the error message.
func RequireNoErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error, got: %+q", err.Error())
	}
}

// RequireErr checks if the error is not nil, and if it is nil, it fails the test with a message.
func RequireErr(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

// RequireEqual checks if the expected and actual values are equal, and if not, it fails the test with a formatted message.
func RequireEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Fatalf("expected: %+v, actual: %+v", expected, actual)
	}
}

// RequireTrue checks if the condition is true, and if not, it fails the test with a message.
func RequireTrue(t *testing.T, condition bool) {
	t.Helper()

	if !condition {
		t.Fatal("expected condition to be true, got false")
	}
}

// RequireContain checks if the haystack string contains the needle string, and if not, it fails the test with a formatted message.
func RequireContain(t *testing.T, haystack, needle string) {
	t.Helper()

	if strings.Contains(haystack, needle) {
		return
	}

	t.Fatalf("expected %q to contain %q", haystack, needle)
}

// RequireNotContain checks if the haystack string does not contain the needle string, and if it does, it fails the test with a formatted message.
func RequireNotContain(t *testing.T, haystack, needle string) {
	t.Helper()

	if !strings.Contains(haystack, needle) {
		return
	}

	t.Fatalf("expected %q to not contain %q", haystack, needle)
}

// RequireNotNil checks if the object is not nil, and if it is nil, it fails the test with a message.
func RequireNotNil(t *testing.T, o any) {
	t.Helper()

	if o == nil {
		t.Fatal("expected object to not be nil, got nil")
	}
}
