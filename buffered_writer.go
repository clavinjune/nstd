package nstd

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"
)

var _ io.Writer = (*BufferedWriter)(nil)

// BufferedWriter is a thread-safe buffered writer that writes to an underlying io.Writer.
// It automatically flushes the buffer either when it reaches a certain size or after a specified time interval.
// It supports cancellation via context, allowing graceful shutdowns.
type BufferedWriter struct {
	sync.Mutex
	ctx               context.Context
	writer            *bufio.Writer
	maxBufferSize     int
	currentBufferSize int
	flushInterval     time.Duration
	lastFlushAt       time.Time
}

// NewBufferedWriter creates a new BufferedWriter with the specified parameters.
// it runs an auto-flush goroutine that flushes the buffer either when it reaches the maxBufferSize
// or after the flushInterval has passed.
func NewBufferedWriter(ctx context.Context, w io.Writer, maxBufferSize int, flushInterval time.Duration) (*BufferedWriter, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	bw := &BufferedWriter{
		ctx:           ctx,
		writer:        bufio.NewWriter(w),
		maxBufferSize: maxBufferSize,
		flushInterval: flushInterval,
	}

	go bw.autoFlush()

	return bw, func() {
		_ = bw.Flush()
		cancel()
	}
}

// Write writes data to the underlying writer.
func (b *BufferedWriter) Write(p []byte) (int, error) {
	if b.ctx.Err() != nil {
		return 0, b.ctx.Err()
	}

	b.Lock()
	defer b.Unlock()
	n, err := b.writer.Write(p)
	if err != nil {
		return n, err
	}
	b.currentBufferSize += n

	if b.currentBufferSize >= b.maxBufferSize {
		err := b.unsafeFlush()
		if err != nil {
			return 0, err
		}
	}

	return n, err
}

func (b *BufferedWriter) unsafeFlush() error {
	if err := b.writer.Flush(); err != nil {
		return err
	}
	b.currentBufferSize = 0
	b.lastFlushAt = time.Now()

	return nil
}

// Flush flushes the buffered data to the underlying writer.
func (b *BufferedWriter) Flush() error {
	b.Lock()
	defer b.Unlock()
	return b.unsafeFlush()
}

// autoFlush runs in a separate goroutine and periodically checks if the buffer should be flushed.
// it flushes the buffer if the time since the last flush exceeds the flushInterval.
// it also flushes the buffer when the context is done.
func (b *BufferedWriter) autoFlush() {
	t := time.NewTicker(b.flushInterval)
	defer t.Stop()

	for {
		select {
		case <-b.ctx.Done():
			_ = b.Flush()
			return
		case <-t.C:
			b.Lock()
			if time.Since(b.lastFlushAt) >= b.flushInterval {
				_ = b.unsafeFlush()
			}
			b.Unlock()
		}
	}
}
