package gouache

import (
	"context"
	"errors"
)

var _ Cache = (*ReadThroughCache)(nil)

type ReadThroughCache struct {
	Cache
	Loader Loader
}

func (cache *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := cache.Cache.Get(ctx, key)
	if errors.Is(err, ErrCacheMiss) {
		val, err := cache.Loader.Load(ctx, key)
		if err != nil {
			return nil, err
		}
		return val, cache.Cache.Set(ctx, key, val)
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}
