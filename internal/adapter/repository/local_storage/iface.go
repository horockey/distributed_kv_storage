package local_storage

type Repository[V any] interface {
	Get(key string) (V, error)
	Set(key string, val V) error
}
