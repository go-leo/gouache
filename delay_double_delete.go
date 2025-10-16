package gouache

import (
	"context"
	"errors"
)

type DelayDoubleDeleteStore struct {
	Cache  Cache
	Loader Loader
	Updater Updater
}

func (store *DelayDoubleDeleteStore) Get(ctx context.Context, key string) (any, error) {
	val, err := store.Cache.Get(ctx, key)
	if errors.Is(err, ErrCacheMiss) {
		val, err := store.Loader.Load(ctx, key)
		if err != nil {
			return nil, err
		}
		return val, store.Cache.Set(ctx, key, val)
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (store *DelayDoubleDeleteStore) Set(ctx context.Context, key string, val any) error {
	if err := store.Cache.Delete(ctx, key); err != nil {
		return err
	}
	store.
	return store.Cache.Set(ctx, key, val)
}
