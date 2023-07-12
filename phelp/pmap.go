package phelp

func Keys[K comparable, V any](mapData map[K]V) []K {
	keys := make([]K, 0, len(mapData))
	for k := range mapData {
		keys = append(keys, k)
	}
	return keys
}
