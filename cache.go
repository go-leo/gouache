package gouache

import (
	"context"
	"errors"
)

var ErrCacheMiss = errors.New("gouache: key not found")

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any) error
	Delete(ctx context.Context, key string) error
}

type Database interface {
	Select(ctx context.Context, key string) (any, error)
	Upsert(ctx context.Context, key string, val any) error
}
