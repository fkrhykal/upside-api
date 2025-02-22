package collection

func Map[T comparable, P any](items []T, fn func(item T) P) []P {
	result := make([]P, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}
