// Package filter provides functions for filtering.
package filter

// Func is a function that filters T.
type Func[T any] func(T) bool

// Inverse returns a filter that inverts the given filter.
func Inverse[T any](filter Func[T]) Func[T] {
	return func(t T) bool {
		return !filter(t)
	}
}
