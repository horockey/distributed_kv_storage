package remote_storage

type AppNode struct {
	Name    string
	Address string
}

type Gateway interface {
	Get(key string, node AppNode) (map[string]any, error)
	Set(key string, val map[string]any, node AppNode) error
}
