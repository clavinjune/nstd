package nstd_test

import (
	"testing"

	. "github.com/clavinjune/nstd"
)

func TestNew(t *testing.T) {
	tt := []struct {
		_     struct{}
		Name  string
		Input any
	}{
		{
			Name:  "bool",
			Input: true,
		},
		{
			Name:  "int",
			Input: 123,
		},
		{
			Name:  "float64",
			Input: 12.34,
		},
		{
			Name:  "string",
			Input: "string",
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			got := New(tc.Input)
			RequireEqual(t, *got, tc.Input)
		})
	}
}
