package local_storage

type Repository interface {
	Get(key string) (map[string]any, error)
	Set(key string, val map[string]any) error
}
