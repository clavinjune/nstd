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
		RequireEqual(t, n, 13)
		RequireNil(t, err)
		RequireEqual(t, b.String(), "")

		time.Sleep(150 * time.Millisecond)
		RequireEqual(t, b.String(), "Hello, World!")
	})

	t.Run("more than buffer size", func(t *testing.T) {
		defer b.Reset()

		n, err := bw.Write([]byte("more than 15bytes should be flushed immediately"))
		RequireEqual(t, n, 47)
		RequireNil(t, err)
		RequireEqual(t, b.String(), "more than 15bytes should be flushed immediately")
	})

	t.Run("flush manually", func(t *testing.T) {
		defer b.Reset()
		n, err := bw.Write([]byte("a"))
		RequireEqual(t, n, 1)
		RequireNil(t, err)
		RequireNil(t, bw.Flush())
		RequireEqual(t, b.String(), "a")
	})

	t.Run("close buffered writer", func(t *testing.T) {
		defer b.Reset()
		n, err := bw.Write([]byte("written"))
		RequireEqual(t, n, 7)
		RequireNil(t, err)
		closeBw()
		RequireEqual(t, b.String(), "written")

		b.Reset()
		n, err = bw.Write([]byte("This should not be written"))
		RequireEqual(t, n, 0)
		RequireNotNil(t, err)
		RequireEqual(t, b.String(), "")
	})
}
