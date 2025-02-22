package collection

func Empty[T comparable](d []T) bool {
	return len(d) == 0
}
