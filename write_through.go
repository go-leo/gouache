package gouache

import (
	"context"
)

type Updater interface {
	Update(ctx context.Context, key string, val any) error
}

var _ Cache = (*WriteThroughCache)(nil)

type WriteThroughCache struct {
	Cache
	Updater Updater
}

func (cache *WriteThroughCache) Set(ctx context.Context, key string, val any) error {
	if err := cache.Cache.Delete(ctx, key); err != nil {
		return err
	}
	if err := cache.Updater.Update(ctx, key, val); err != nil {
		return err
	}
	return cache.Cache.Set(ctx, key, val)
}
