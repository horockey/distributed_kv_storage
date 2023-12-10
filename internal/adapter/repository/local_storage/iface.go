package local_storage

type Repository[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, val V) error
}
