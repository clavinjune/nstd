package nstd_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/clavinjune/nstd"
)

func ExampleBufferedWriter() {
	b := new(bytes.Buffer) // you can use any io.Writer here including os.Stdout

	bw, closeBw := nstd.NewBufferedWriter(context.Background(), b, 15, 100*time.Millisecond)
	defer closeBw()

	bw.Write([]byte("Hello, World!"))
	time.Sleep(150 * time.Millisecond) // Wait for auto-flush
	fmt.Println(b.String())
	// Output: Hello, World!
}

func ExampleFlagSet() {
	defer os.Clearenv()
	os.Setenv("EXAMPLE_NAME", "from-envs")
	fs := nstd.NewFlagSet("example", flag.ExitOnError)
	nameFlag := fs.String("name", "default", "usage")

	if err := fs.Parse("--name", "from-args"); err != nil {
		panic(err)
	}

	fmt.Println(*nameFlag)
	// Output: from-envs
}

func ExampleNewShutdownContext() {
	// ShutdownContext will be canceled when a SIGINT or SIGTERM is received or manually canceled.
	ctx, cancel := nstd.NewShutdownContext(context.Background())
	defer cancel()
	errChan := make(chan error, 1)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	// nstd.Wait will wait until the context is done or an error is received from the error channel.
	if err := nstd.Wait(ctx, errChan); err != nil {
		// handle the error, e.g., log it
	}
}

func ExampleNewShutdownContextWithCause() {
	// ShutdownContextWithCause will be more verbose about the cause of the shutdown.
	ctx, cancel := nstd.NewShutdownContextWithCause(context.Background())
	defer cancel(context.Canceled)

	errChan := make(chan error, 1)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	// nstd.Wait will wait until the context is done or an error is received from the error channel.
	if err := ctx.Wait(errChan); err != nil {
		// handle the error, e.g., log it
	}

}

func ExampleNewSlog() {
	isDebug := true
	isStructured := true
	writer := os.Stdout
	logger := nstd.NewSlog(writer, isDebug, isStructured)
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
}
