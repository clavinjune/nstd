package nstd_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
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

func ExampleNewSlog() {
	isDebug := true
	isStructured := true
	writer := os.Stdout
	logger := nstd.NewSlog(writer, isDebug, isStructured)
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
}
