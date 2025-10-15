package lrux

import (
	"context"

	"github.com/go-leo/gouache"
	lru "github.com/hashicorp/golang-lru"
)

var _ gouache.Cache = (*Cache)(nil)

// Cache LRU缓存
type Cache struct {
	LRUCache *lru.Cache
}

func (store *Cache) Get(ctx context.Context, key string) (any, error) {
	val, ok := store.LRUCache.Get(key)
	if !ok {
		return nil, gouache.ErrNil
	}
	return val, nil
}

func (store *Cache) Set(ctx context.Context, key string, val any) error {
	_ = store.LRUCache.Add(key, val)
	return nil
}

func (store *Cache) Delete(ctx context.Context, key string) error {
	_ = store.LRUCache.Remove(key)
	return nil
}
