package nstd_test

import (
	"testing"

	. "github.com/clavinjune/nstd"
)

func TestPoolBasic(t *testing.T) {
	pool := NewPool(func() int { return 1 })
	a := pool.Get()
	if a != 1 {
		t.Fatalf("expected to be getting: 1 but got: %v", a)
	}
	pool.Put(a)
	b := pool.Get()
	if b != a {
		t.Fatalf("expected b and a to be: 1 but got: %v", b)
	}
}
func TestNewPool(t *testing.T) {
	type testCase[T any] struct {
		Name   string
		Fn     func() T
		Action func(got T) bool
	}

	tests := []testCase[any]{
		{
			Name: "String",
			Fn: func() any {
				return "new pool"
			},
			Action: func(got any) bool {
				v, ok := got.(string)
				return ok && v == "new pool"
			},
		},
		{
			Name: "int",
			Fn: func() any {
				return 0
			},
			Action: func(got any) bool {
				v, ok := got.(int)
				return ok && v == 0
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			pool := NewPool(tt.Fn)
			got := pool.Get()
			if !tt.Action(got) {
				t.Fatalf("unexpected value got from pool: %v", got)
			}

		})
	}
}

func TestPtrPool(t *testing.T) {
	type test struct {
		a int
	}
	pool := NewPtrPool[test]()
	x := pool.Get()
	if x == nil {
		t.Fatalf("expected non-nil pointer from Get()")
	}

	if x.a != 0 {
		t.Fatalf("expected to be 0 value struct but got: %v", *x)
	}
	x.a = 1
	pool.Put(x)

	y := pool.Get()
	if y == nil {
		t.Fatalf("expected non-nil pointer after reusing")
	}

	if y.a != 1 {
		t.Fatalf("expecting to be reused but got: %v", *y)
	}
}
