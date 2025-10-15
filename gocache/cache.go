package gocache

import (
	"context"
	"time"

	"github.com/go-leo/gouache"
	"github.com/patrickmn/go-cache"
)

var _ gouache.Cache = (*Cache)(nil)

type Cache struct {
	Cache *cache.Cache
	TTL   func(ctx context.Context, key string, val any) (time.Duration, error)
}

func (store *Cache) Get(ctx context.Context, key string) (any, error) {
	val, ok := store.Cache.Get(key)
	if !ok {
		return nil, gouache.ErrNil
	}
	return val, nil
}

func (store *Cache) Set(ctx context.Context, key string, val any) error {
	if store.TTL != nil {
		ttl, err := store.TTL(ctx, key, val)
		if err != nil {
			return err
		}
		store.Cache.Set(key, val, ttl)
	}
	store.Cache.Set(key, val, cache.DefaultExpiration)
	return nil
}

func (store *Cache) Delete(ctx context.Context, key string) error {
	store.Cache.Delete(key)
	return nil
}
