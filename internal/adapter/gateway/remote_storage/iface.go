package remote_storage

type AppNode struct {
	Name    string
	Address string
}

type Gateway[K comparable, V any] interface {
	Get(key K, node AppNode) (V, error)
	Set(key K, val V, node AppNode) error
}
