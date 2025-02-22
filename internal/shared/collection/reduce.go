package collection

func Reduce[T comparable, P any](items []T, fn func(item T, acc P), acc P) P {
	for _, item := range items {
		fn(item, acc)
	}
	return acc
}
