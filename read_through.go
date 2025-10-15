package gouache

import (
	"context"
	"errors"
)

type Loader interface {
	Load(ctx context.Context, key string) (any, error)
}

var _ Cache = (*ReadThroughCache)(nil)

type ReadThroughCache struct {
	Cache
	Loader Loader
}

func (cache *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := cache.Cache.Get(ctx, key)
	if errors.Is(err, ErrNil) {
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
