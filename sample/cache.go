package sample

import (
	"context"
	"sync"

	"github.com/go-leo/gouache"
)

var _ gouache.Cache = (*Cache)(nil)

// Cache 简单缓存
type Cache struct {
	Map sync.Map
}

func (store *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	val, ok := store.Map.Load(key)
	if !ok {
		return nil, gouache.ErrNil
	}
	return val, nil
}

func (store *Cache) Set(ctx context.Context, key string, val any) error {
	store.Map.Store(key, val)
	return nil
}

func (store *Cache) Delete(ctx context.Context, key string) error {
	store.Map.Delete(key)
	return nil
}
