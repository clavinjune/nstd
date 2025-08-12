package nstd_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/clavinjune/nstd"
)

func ExampleBufferedWriter() {
	var b bytes.Buffer
	w := bufio.NewWriter(&b) // you can use any io.Writer here including os.Stdout

	bw, closeBw := NewBufferedWriter(context.Background(), w, 15, 100*time.Millisecond)
	defer closeBw()

	bw.Write([]byte("Hello, World!"))
	time.Sleep(150 * time.Millisecond) // Wait for auto-flush
	fmt.Println(b.String())
	// Output: Hello, World!
}

func TestBufferedWriter(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	bw, closeBw := NewBufferedWriter(context.Background(), w, 15, 100*time.Millisecond)
	defer closeBw()

	t.Run("more than flush interval", func(t *testing.T) {
		defer b.Reset()

		n, err := bw.Write([]byte("Hello, World!"))
		RequireEqual(t, 13, n)
		RequireNoErr(t, err)
		RequireEqual(t, "", b.String())

		time.Sleep(150 * time.Millisecond)
		RequireEqual(t, "Hello, World!", b.String())
	})

	t.Run("more than buffer size", func(t *testing.T) {
		defer b.Reset()

		n, err := bw.Write([]byte("more than 15bytes should be flushed immediately"))
		RequireEqual(t, 47, n)
		RequireNoErr(t, err)
		RequireEqual(t, "more than 15bytes should be flushed immediately", b.String())
	})

	t.Run("flush manually", func(t *testing.T) {
		defer b.Reset()
		n, err := bw.Write([]byte("a"))
		RequireEqual(t, 1, n)
		RequireNoErr(t, err)
		RequireNoErr(t, bw.Flush())
		RequireEqual(t, "a", b.String())
	})

	t.Run("close buffered writer", func(t *testing.T) {
		defer b.Reset()
		closeBw()

		n, err := bw.Write([]byte("This should not be written"))
		RequireEqual(t, 0, n)
		RequireErr(t, err)
		RequireEqual(t, "", b.String())
	})
}
