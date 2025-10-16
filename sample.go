package gouache

import (
	"context"
	"sync"
)

var _ Cache = (*SampleCache)(nil)

type SampleCache struct {
	Cache sync.Map
}

func (store *SampleCache) Get(ctx context.Context, key string) (any, error) {
	val, ok := store.Cache.Load(key)
	if !ok {
		return nil, ErrCacheMiss
	}
	return val, nil
}

func (store *SampleCache) Set(ctx context.Context, key string, val any) error {
	store.Cache.Store(key, val)
	return nil
}

func (store *SampleCache) Delete(ctx context.Context, key string) error {
	store.Cache.Delete(key)
	return nil
}
