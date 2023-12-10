package inmemory_local_storage

import (
	"sync"

	"github.com/horockey/distributed_kv_storage/internal/adapter/repository/local_storage"
)

var _ local_storage.Repository = &inmemoryLocalStorage{}

type inmemoryLocalStorage struct {
	mu      sync.RWMutex
	storage map[string]map[string]any
}

func New() *inmemoryLocalStorage {
	return &inmemoryLocalStorage{
		storage: map[string]map[string]any{},
	}
}

func (repo *inmemoryLocalStorage) Get(key string) (map[string]any, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.storage[key], nil
}

func (repo *inmemoryLocalStorage) Set(key string, val map[string]any) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.storage[key] = val
	return nil
}
