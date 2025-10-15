package bigcache

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/allegro/bigcache"
	"github.com/go-leo/gouache"
)

var _ gouache.Cache = (*Cache)(nil)

type Cache struct {
	BigCache  *bigcache.BigCache
	Marshal   func(key string, obj any) ([]byte, error)
	Unmarshal func(key string, data []byte) (any, error)
}

func (store *Cache) Get(ctx context.Context, key string) (any, error) {
	data, err := store.BigCache.Get(key)
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		return nil, gouache.ErrNil
	}
	if err != nil {
		return nil, err
	}
	if store.Unmarshal == nil {
		return data, nil
	}
	obj, err := store.Unmarshal(key, data)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (store *Cache) Set(ctx context.Context, key string, val any) error {
	if store.Marshal == nil {
		switch val := val.(type) {
		case []byte:
			return store.BigCache.Set(key, val)
		case string:
			return store.BigCache.Set(key, []byte(val))

		case int:
			return store.BigCache.Set(key, []byte(strconv.FormatInt(int64(val), 10)))
		case int8:
			return store.BigCache.Set(key, []byte(strconv.FormatInt(int64(val), 10)))
		case int16:
			return store.BigCache.Set(key, []byte(strconv.FormatInt(int64(val), 10)))
		case int32:
			return store.BigCache.Set(key, []byte(strconv.FormatInt(int64(val), 10)))
		case int64:
			return store.BigCache.Set(key, []byte(strconv.FormatInt(int64(val), 10)))

		case uint:
			return store.BigCache.Set(key, []byte(strconv.FormatUint(uint64(val), 10)))
		case uint8:
			return store.BigCache.Set(key, []byte(strconv.FormatUint(uint64(val), 10)))
		case uint16:
			return store.BigCache.Set(key, []byte(strconv.FormatUint(uint64(val), 10)))
		case uint32:
			return store.BigCache.Set(key, []byte(strconv.FormatUint(uint64(val), 10)))
		case uint64:
			return store.BigCache.Set(key, []byte(strconv.FormatUint(uint64(val), 10)))

		case float32:
			return store.BigCache.Set(key, []byte(strconv.FormatFloat(float64(val), 'f', -1, 32)))
		case float64:
			return store.BigCache.Set(key, []byte(strconv.FormatFloat(float64(val), 'f', -1, 64)))

		case bool:
			return store.BigCache.Set(key, []byte(strconv.FormatBool(val)))
		default:
			return fmt.Errorf("bigcache: failed to convert %v(%t) to bytes", val, val)
		}
	}
	data, err := store.Marshal(key, val)
	if err != nil {
		return err
	}
	return store.BigCache.Set(key, data)
}

func (store *Cache) Delete(ctx context.Context, key string) error {
	return store.BigCache.Delete(key)
}
