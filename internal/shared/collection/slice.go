package collection

func Slice[T comparable](v []T, i, j int) []T {
	var null T
	if At(v, i) == null {
		return v
	}
	if At(v, j) == null {
		return v[i:]
	}
	return v[i:j]
}
