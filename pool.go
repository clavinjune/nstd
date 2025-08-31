package nstd

import "sync"

// Pool wraps sync.Pool using generic
type Pool[T any] struct {
	sync.Pool
}

// Pool gets object from pool and parses it into T
func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

// Put wraps sync.Pool.Put
func (p *Pool[T]) Put(t T) {
	p.Pool.Put(t)
}

// NewPool returns Pool with given newFn
func NewPool[T any](newFn func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() any {
				return newFn()
			},
		},
	}
}
