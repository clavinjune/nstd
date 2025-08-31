package nstd_test

import (
	"fmt"
	"testing"

	. "github.com/clavinjune/nstd"
)

type wrapper struct {
	testing.TB
	errorMap map[string]struct{}
}

func (w *wrapper) Errorf(format string, args ...any) {
	w.errorMap[fmt.Sprintf(format, args...)] = struct{}{}
}

func (w *wrapper) Error(args ...any) {
	w.errorMap[fmt.Sprintf("%+v", args...)] = struct{}{}
}

func TestRequire(t *testing.T) {
	m := make(map[string]struct{})
	w := &wrapper{t, m}

	t.Run("RequireTrue", func(t *testing.T) {
		RequireTrue(w, false)
		_, ok := m["got: false, want: true"]
		RequireTrue(t, ok)
	})
	t.Run("RequireNotNil", func(t *testing.T) {
		RequireNotNil(w, nil)
		_, ok := m["got: nil, want: not nil"]
		RequireTrue(t, ok)
	})
	t.Run("RequireNil", func(t *testing.T) {
		RequireNil(w, 1)
		_, ok := m["got: 1, want: nil"]
		RequireTrue(t, ok)
	})
	t.Run("RequireEqual", func(t *testing.T) {
		RequireEqual(w, 1, 2)
		_, ok := m["got: 1, want: 2"]
		RequireTrue(t, ok)
	})
}
