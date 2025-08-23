package nstd_test

import (
	"fmt"
	"math"
	"testing"
	"testing/synctest"
	"time"

	"github.com/clavinjune/nstd"
)

func TestPipeline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// buffered channel
		ch := nstd.PipeFrom(2, 1, 2, 3, 4, 5, 6)
		nstd.RequireEqual(t, len(ch), 0)
		nstd.RequireEqual(t, cap(ch), 2)
		time.Sleep(time.Second)
		nstd.RequireEqual(t, len(ch), 2)

		doubleCh := nstd.PipeMap(ch, func(i int) float64 {
			return float64(i * 2)
		})
		// 2, 4, 6, 8, 10, 12

		mod4Ch := nstd.PipeFilter(doubleCh, func(f float64) bool {
			return math.Mod(f, 4) == 0
		})
		// 4, 8, 12

		result := nstd.PipeReduce(mod4Ch, []float64{}, func(cumulative []float64, currentValue float64) []float64 {
			return append(cumulative, currentValue)
		})
		// [4, 8, 12]

		// unbuffered channel
		arrCh := nstd.PipeFrom(0, result...)
		nstd.RequireEqual(t, len(arrCh), 0)
		nstd.RequireEqual(t, cap(arrCh), 0)

		strCh := nstd.PipeMap(arrCh, func(i float64) string {
			return fmt.Sprintf("%0.1f", i)
		})
		// 4.0, 8.0, 12.0
		result2 := nstd.PipeReduce(strCh, "", func(cumulative string, current string) string {
			if cumulative == "" {
				return current
			}
			return fmt.Sprintf("%s,%s", cumulative, current)
		})

		nstd.RequireEqual(t, result2, "4.0,8.0,12.0")
	})
}

func TestPipeTo(t *testing.T) {
	ch := nstd.PipeFrom(0, 0, 1, 2, 3, 4)
	for i, v := range nstd.PipeTo(ch) {
		nstd.RequireEqual(t, v, i)
	}
}
