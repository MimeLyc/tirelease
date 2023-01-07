package utils

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func Intersects[T comparable](s []T, m []T) []T {
	if len(s) == 0 || len(m) == 0 {
		return []T{}
	}

	result := make([]T, 0)

	for _, melement := range m {
		if Contains(s, melement) {
			result = append(result, melement)
		}
	}

	return result
}
