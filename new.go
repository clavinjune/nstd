package nstd

// New simplifies the new function with default value
func New[T any](t T) *T {
	return &t
}
