package gouache

import (
	"context"

	"golang.org/x/sync/singleflight"
)

var _ Cache = (*SingleFlightCache)(nil)

type SingleFlightCache struct {
	Cache       Cache
	getGroup    singleflight.Group
	setGroup    singleflight.Group
	deleteGroup singleflight.Group
}

func (store *SingleFlightCache) Get(ctx context.Context, key string) (any, error) {
	val, err, _ := store.getGroup.Do(key, func() (any, error) {
		return store.Cache.Get(ctx, key)
	})
	return val, err
}

func (store *SingleFlightCache) Set(ctx context.Context, key string, val any) error {
	_, err, _ := store.setGroup.Do(key, func() (any, error) {
		return nil, store.Cache.Set(ctx, key, val)
	})
	return err
}

func (store *SingleFlightCache) Delete(ctx context.Context, key string) error {
	_, err, _ := store.deleteGroup.Do(key, func() (any, error) {
		return nil, store.Cache.Delete(ctx, key)
	})
	return err
}
