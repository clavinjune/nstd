package nstd_test

import (
	"bytes"
	"encoding/json"
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
		RequireNoErr(t, json.Unmarshal(b.Bytes(), &m))
		RequireEqual(t, "should be shown", m["msg"])
		RequireEqual(t, "DEBUG", m["level"])
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
		RequireNoErr(t, json.Unmarshal(b.Bytes(), &m))
		RequireEqual(t, "should be shown", m["msg"])
		RequireEqual(t, "INFO", m["level"])
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

		RequireErr(t, json.Unmarshal(b.Bytes(), &struct{}{}))
		str := b.String()
		RequireContain(t, str, `level=INFO`)
		RequireContain(t, str, `msg="should be shown"`)
		RequireNotContain(t, str, "source=")
	})
	t.Run("debug level with unstructured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := NewSlog(b, true, false)
		logger.Debug("should be shown")

		RequireErr(t, json.Unmarshal(b.Bytes(), &struct{}{}))
		str := b.String()
		RequireContain(t, str, `level=DEBUG`)
		RequireContain(t, str, `msg="should be shown"`)
		RequireContain(t, str, `source=`)
	})
}
