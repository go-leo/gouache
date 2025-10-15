package gouache

import (
	"context"
	"errors"
)

// ErrNil reply returned by cache when key does not exist.
var ErrNil = errors.New("cache: nil")

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any) error
	Delete(ctx context.Context, key string) error
}
