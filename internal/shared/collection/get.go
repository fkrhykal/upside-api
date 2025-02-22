package collection

func At[T comparable](v []T, i int) T {
	var null T
	if LastIndex(v) < 0 {
		return null
	}
	if LastIndex(v) < i {
		return null
	}
	return v[i]
}

func First[T comparable](v []T) T {
	return At(v, 0)
}

func Second[T comparable](v []T) T {
	return At(v, 1)
}

func Last[T comparable](v []T) T {
	return At(v, LastIndex(v))
}

func Penultimate[T comparable](v []T) T {
	return At(v, LastIndex(v)-1)
}

func LastIndex[T comparable](v []T) int {
	return len(v) - 1
}
