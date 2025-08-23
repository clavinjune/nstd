package nstd_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	. "github.com/clavinjune/nstd"
)

func TestBufferedWriter(t *testing.T) {
	b := new(bytes.Buffer)
	bw, closeBw := NewBufferedWriter(context.Background(), b, 15, 100*time.Millisecond)
	defer closeBw()

	t.Run("more than flush interval", func(t *testing.T) {
		defer b.Reset()

		n, err := bw.Write([]byte("Hello, World!"))
		RequireEqual(t, 13, n)
		RequireNil(t, err)
		RequireEqual(t, "", b.String())

		time.Sleep(150 * time.Millisecond)
		RequireEqual(t, "Hello, World!", b.String())
	})

	t.Run("more than buffer size", func(t *testing.T) {
		defer b.Reset()

		n, err := bw.Write([]byte("more than 15bytes should be flushed immediately"))
		RequireEqual(t, 47, n)
		RequireNil(t, err)
		RequireEqual(t, "more than 15bytes should be flushed immediately", b.String())
	})

	t.Run("flush manually", func(t *testing.T) {
		defer b.Reset()
		n, err := bw.Write([]byte("a"))
		RequireEqual(t, 1, n)
		RequireNil(t, err)
		RequireNil(t, bw.Flush())
		RequireEqual(t, "a", b.String())
	})

	t.Run("close buffered writer", func(t *testing.T) {
		defer b.Reset()
		closeBw()

		n, err := bw.Write([]byte("This should not be written"))
		RequireEqual(t, 0, n)
		RequireNotNil(t, err)
		RequireEqual(t, "", b.String())
	})
}
