package utils

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func SliceFromMap[M ~map[K]V, K comparable, V any](m M) []V {

	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func GetElement[T any](s []T, i int) T {

	if i >= 0 {
		return s[i]
	}

	i = len(s) + i
	return s[i]

}

func SaveElement[T any](s []T, i int, v T) {

	if i >= 0 {
		s[i] = v
		return
	}

	i = len(s) + i
	s[i] = v

}
