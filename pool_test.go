package nstd_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/clavinjune/nstd"
)

func TestPool(t *testing.T) {
	p := nstd.NewPool(func() *bytes.Buffer {
		return new(bytes.Buffer)
	})

	b := p.Get()
	fmt.Fprint(b, "Hello, World!")

	nstd.RequireEqual(t, b.String(), "Hello, World!")

	b.Reset()
	p.Put(b)

	bb := p.Get()
	fmt.Fprint(bb, "Hello, World!")
	nstd.RequireEqual(t, bb.String(), "Hello, World!")
	nstd.RequireEqual(t, b, bb)
}
