package gouache

import (
	"context"
	"encoding/binary"
	"hash"
	"hash/fnv"
)

var _ Cache = (*shardedCache)(nil)

type HashFactory func(ctx context.Context, key string) (hash.Hash, error)

type shardedCache struct {
	Options *shardedOptions
	Buckets []Cache
}

type shardedOptions struct {
	HashFactory HashFactory
}

type ShardedOption func(*shardedOptions)

func WithHashFactory(hashFactory HashFactory) ShardedOption {
	return func(o *shardedOptions) {
		o.HashFactory = hashFactory
	}
}

func newShardedOptions(opts ...ShardedOption) *shardedOptions {
	options := &shardedOptions{}
	return options.Apply(opts...).Correct()
}

func (o *shardedOptions) Apply(opts ...ShardedOption) *shardedOptions {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *shardedOptions) Correct() *shardedOptions {
	if o.HashFactory == nil {
		o.HashFactory = func(ctx context.Context, key string) (hash.Hash, error) {
			return fnv.New32a(), nil
		}
	}
	return o
}

func ShardedCache(buckets []Cache, opts ...ShardedOption) Cache {
	if len(buckets) == 0 {
		panic("gouache: buckets is empty")
	}
	return &shardedCache{Options: newShardedOptions(opts...), Buckets: buckets}
}

func (store *shardedCache) Get(ctx context.Context, key string) (any, error) {
	return store.bucket(ctx, key).Get(ctx, key)
}

func (store *shardedCache) Set(ctx context.Context, key string, val any) error {
	return store.bucket(ctx, key).Set(ctx, key, val)
}

func (store *shardedCache) Delete(ctx context.Context, key string) error {
	return store.bucket(ctx, key).Delete(ctx, key)
}

func (cache *shardedCache) bucket(ctx context.Context, key string) Cache {
	h, err := cache.Options.HashFactory(ctx, key)
	if err != nil {
		return nil
	}
	h.Write([]byte(key))
	switch h.Size() {
	case 4:
		sum32 := int(h.(hash.Hash32).Sum32())
		return cache.Buckets[sum32%len(cache.Buckets)]
	case 8:
		sum64 := int(h.(hash.Hash64).Sum64())
		return cache.Buckets[sum64%len(cache.Buckets)]
	default:
		sum := h.Sum(nil)
		if len(sum) < 4 {
			return cache.Buckets[0]
		}
		sum32 := int(binary.BigEndian.Uint32(sum[:4]))
		return cache.Buckets[sum32%len(cache.Buckets)]
	}
}
