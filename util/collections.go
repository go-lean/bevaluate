package util

func MapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))

	for k := range m {
		result = append(result, k)
	}

	return result
}
