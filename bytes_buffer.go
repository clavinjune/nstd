package nstd

import (
	"bytes"
	"sync"
)

// BytesBuffer is a safe-thread bytes.Buffer
// only use this in test
type BytesBuffer struct {
	sync.Mutex
	b bytes.Buffer
}

func (b *BytesBuffer) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	return b.b.Write(p)
}
func (b *BytesBuffer) Reset() {
	b.Lock()
	defer b.Unlock()
	b.b.Reset()
}

func (b *BytesBuffer) String() string {
	b.Lock()
	defer b.Unlock()
	return b.b.String()
}
