package util

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~string
}

func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func IndexOf[T comparable](slice []T, item T) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

func Mapping[T, U any](s []T, f func(T) U) []U {
	ret := make([]U, len(s))
	for k, v := range s {
		ret[k] = f(v)
	}
	return ret
}

type Map[K comparable, V any] struct {
	data map[K]V
}

func (m *Map[K, V]) Get(k K) V {
	return m.data[k]
}
func (m *Map[K, V]) Set(k K, v V) {
	m.data[k] = v
}
func (m *Map[K, V]) Delete(k K) {
	delete(m.data, k)
}
