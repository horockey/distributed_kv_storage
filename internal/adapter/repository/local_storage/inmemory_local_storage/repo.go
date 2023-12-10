package inmemory_local_storage

import (
	"sync"

	"github.com/horockey/distributed_kv_storage/internal/adapter/repository/local_storage"
)

var _ local_storage.Repository[struct{}] = &inmemoryLocalStorage[struct{}]{}

type inmemoryLocalStorage[V any] struct {
	mu      sync.RWMutex
	storage map[string]V
}

func New[V any]() *inmemoryLocalStorage[V] {
	return &inmemoryLocalStorage[V]{
		storage: map[string]V{},
	}
}

func (repo *inmemoryLocalStorage[V]) Get(key string) (V, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.storage[key], nil
}

func (repo *inmemoryLocalStorage[V]) Set(key string, val V) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.storage[key] = val
	return nil
}
