package util

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func MapValueOrDefault[K comparable, V any](m map[K]V, key K, defval V) V {
	v, ok := m[key]

	if !ok {
		v = defval
		m[key] = v
	}

	return v
}

func MapSortKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := []K{}

	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}
