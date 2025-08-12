package nstd_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/clavinjune/nstd"
)

func ExampleNewSlog() {
	logger := nstd.NewSlog(os.Stdout, true, true)
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
}

func TestNewSlog(t *testing.T) {
	t.Run("debug level with structured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := nstd.NewSlog(b, true, true)
		logger.Debug("should be shown")

		m := make(map[string]any)
		nstd.RequireNoErr(t, json.Unmarshal(b.Bytes(), &m))
		nstd.RequireEqual(t, "should be shown", m["msg"])
		nstd.RequireEqual(t, "DEBUG", m["level"])
		_, ok := m["source"]
		nstd.RequireTrue(t, ok)
	})

	t.Run("info level with structured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := nstd.NewSlog(b, false, true)
		logger.Debug("shouldn't be shown")
		logger.Info("should be shown")

		m := make(map[string]any)
		nstd.RequireNoErr(t, json.Unmarshal(b.Bytes(), &m))
		nstd.RequireEqual(t, "should be shown", m["msg"])
		nstd.RequireEqual(t, "INFO", m["level"])
		_, ok := m["source"]
		nstd.RequireTrue(t, !ok)
	})
	t.Run("info level with unstructured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := nstd.NewSlog(b, false, false)
		logger.Debug("shouldn't be shown")
		logger.Info("should be shown")

		nstd.RequireErr(t, json.Unmarshal(b.Bytes(), &struct{}{}))
		str := b.String()
		nstd.RequireContain(t, str, `level=INFO`)
		nstd.RequireContain(t, str, `msg="should be shown"`)
		nstd.RequireNotContain(t, str, "source=")
	})
	t.Run("debug level with unstructured logging", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)
		defer b.Reset()
		logger := nstd.NewSlog(b, true, false)
		logger.Debug("should be shown")

		nstd.RequireErr(t, json.Unmarshal(b.Bytes(), &struct{}{}))
		str := b.String()
		nstd.RequireContain(t, str, `level=DEBUG`)
		nstd.RequireContain(t, str, `msg="should be shown"`)
		nstd.RequireContain(t, str, `source=`)
	})
}
