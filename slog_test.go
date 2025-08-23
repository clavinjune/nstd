package nstd_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	. "github.com/clavinjune/nstd"
)

func TestNewSlog(t *testing.T) {
	t.Run("debug level with structured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := NewSlog(b, true, true)
		logger.Debug("should be shown")

		m := make(map[string]any)
		RequireNil(t, json.Unmarshal(b.Bytes(), &m))
		RequireEqual(t, m["msg"], "should be shown")
		RequireEqual(t, m["level"], "DEBUG")
		_, ok := m["source"]
		RequireTrue(t, ok)
	})

	t.Run("info level with structured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := NewSlog(b, false, true)
		logger.Debug("shouldn't be shown")
		logger.Info("should be shown")

		m := make(map[string]any)
		RequireNil(t, json.Unmarshal(b.Bytes(), &m))
		RequireEqual(t, m["msg"], "should be shown")
		RequireEqual(t, m["level"], "INFO")
		_, ok := m["source"]
		RequireTrue(t, !ok)
	})
	t.Run("info level with unstructured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := NewSlog(b, false, false)
		logger.Debug("shouldn't be shown")
		logger.Info("should be shown")

		RequireNotNil(t, json.Unmarshal(b.Bytes(), &struct{}{}))
		str := b.String()
		RequireTrue(t, strings.Contains(str, `level=INFO`))
		RequireTrue(t, strings.Contains(str, `msg="should be shown"`))
		RequireTrue(t, !strings.Contains(str, `source=`))
	})
	t.Run("debug level with unstructured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := NewSlog(b, true, false)
		logger.Debug("should be shown")

		RequireNotNil(t, json.Unmarshal(b.Bytes(), &struct{}{}))
		str := b.String()
		RequireTrue(t, strings.Contains(str, `level=DEBUG`))
		RequireTrue(t, strings.Contains(str, `msg="should be shown"`))
		RequireTrue(t, strings.Contains(str, `source=`))
	})
}
