package nstd

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"
)

var _ io.Writer = (*BufferedWriter)(nil)

type BufferedWriter struct {
	sync.Mutex
	ctx               context.Context
	writer            *bufio.Writer
	maxBufferSize     int
	currentBufferSize int
	flushInterval     time.Duration
	lastFlushAt       time.Time
}

func NewBufferedWriter(ctx context.Context, w io.Writer, maxBufferSize int, flushInterval time.Duration) (*BufferedWriter, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	bw := &BufferedWriter{
		ctx:           ctx,
		writer:        bufio.NewWriter(w),
		maxBufferSize: maxBufferSize,
		flushInterval: flushInterval,
	}

	go bw.autoFlush()

	return bw, cancel
}

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
		err := b.Flush()
		if err != nil {
			return 0, err
		}
	}

	return n, err
}

func (b *BufferedWriter) Flush() error {
	if err := b.writer.Flush(); err != nil {
		return err
	}
	b.currentBufferSize = 0
	b.lastFlushAt = time.Now()

	return nil
}

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
				_ = b.Flush()
			}
			b.Unlock()
		}
	}
}
