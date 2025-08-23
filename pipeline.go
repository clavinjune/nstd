package nstd

// PipeFrom transform t to channel
func PipeFrom[T any](buffer int, t ...T) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		for _, i := range t {
			out <- i
		}
	}()

	return out
}

// PipeMap transforms T into U using fn
func PipeMap[T, U any](in <-chan T, fn func(T) U) <-chan U {
	out := make(chan U)
	go func() {
		defer close(out)
		for i := range in {
			out <- fn(i)
		}
	}()

	return out
}

// PipeFilter filters T using given fn
func PipeFilter[T any](in <-chan T, fn func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := range in {
			if fn(i) {
				out <- i
			}
		}
	}()
	return out
}

// PipeReduce reduces T into U using given fn
func PipeReduce[T, U any](in <-chan T, initialValue U, fn func(U, T) U) U {
	out := initialValue
	for i := range in {
		out = fn(out, i)
	}
	return out
}

// PipeTo transforms in to slice
func PipeTo[T any](in <-chan T) []T {
	out := make([]T, 0)
	for i := range in {
		out = append(out, i)
	}

	return out
}
