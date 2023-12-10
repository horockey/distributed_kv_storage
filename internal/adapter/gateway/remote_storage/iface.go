package remote_storage

type AppNode struct {
	Name    string
	Address string
}

type Gateway[V any] interface {
	Get(key string, node AppNode) (V, error)
	Set(key string, val V, node AppNode) error
}
