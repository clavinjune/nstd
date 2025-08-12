package nstd

import (
	"strings"
	"testing"
)

func RequireNoErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error, got: %+q", err.Error())
	}
}

func RequireErr(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func RequireEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Fatalf("expected: %+v, actual: %+v", expected, actual)
	}
}

func RequireTrue(t *testing.T, condition bool) {
	t.Helper()

	if !condition {
		t.Fatal("expected condition to be true, got false")
	}
}

func RequireContain(t *testing.T, haystack, needle string) {
	t.Helper()

	if strings.Contains(haystack, needle) {
		return
	}

	t.Fatalf("expected %q to contain %q", haystack, needle)
}
func RequireNotContain(t *testing.T, haystack, needle string) {
	t.Helper()

	if !strings.Contains(haystack, needle) {
		return
	}

	t.Fatalf("expected %q to not contain %q", haystack, needle)
}
