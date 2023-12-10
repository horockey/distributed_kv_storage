package inmemory_local_storage

import (
	"sync"

	"github.com/horockey/distributed_kv_storage/internal/adapter/repository/local_storage"
)

var _ local_storage.Repository[string, struct{}] = &inmemoryLocalStorage[string, struct{}]{}

type inmemoryLocalStorage[K comparable, V any] struct {
	mu      sync.RWMutex
	storage map[K]V
}

func New[K comparable, V any]() *inmemoryLocalStorage[K, V] {
	return &inmemoryLocalStorage[K, V]{
		storage: map[K]V{},
	}
}

func (repo *inmemoryLocalStorage[K, V]) Get(key K) (V, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.storage[key], nil
}

func (repo *inmemoryLocalStorage[K, V]) Set(key K, val V) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.storage[key] = val
	return nil
}
