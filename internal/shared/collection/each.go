package collection

func Each[T any](items []T, fn func(i int, item T)) {
	for i, item := range items {
		fn(i, item)
	}
}
