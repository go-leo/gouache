package ddd

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/go-leo/gouache"
)

var _ gouache.Cache = (*cache)(nil)

type Gopher func(f func()) error

type options struct {
	DelayDuration time.Duration
	DeleteTimeout time.Duration
	ErrorHandler  func(error)
	Gopher        Gopher
}

type Option func(*options)

func WithDelayDuration(dur time.Duration) Option {
	return func(o *options) {
		o.DelayDuration = dur
	}
}

func WithDeleteTimeout(dur time.Duration) Option {
	return func(o *options) {
		o.DeleteTimeout = dur
	}
}

func WithErrorHandler(f func(error)) Option {
	return func(o *options) {
		o.ErrorHandler = f
	}
}

func WithGopher(gopher Gopher) Option {
	return func(o *options) {
		o.Gopher = gopher
	}
}

func newOptions(opts ...Option) *options {
	options := &options{}
	return options.Apply(opts...).Correct()
}

func (o *options) Apply(opts ...Option) *options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *options) Correct() *options {
	if o.DelayDuration <= 0 {
		o.DelayDuration = 500 * time.Millisecond
	}
	if o.DeleteTimeout <= 0 {
		o.DeleteTimeout = 500 * time.Second
	}
	if o.ErrorHandler == nil {
		o.ErrorHandler = func(err error) {
			slog.Error("ddd.Cache.Get", slog.String("err", err.Error()))
		}
	}
	if o.Gopher == nil {
		o.Gopher = func(f func()) error {
			go f()
			return nil
		}
	}
	return o
}

type cache struct {
	Options  *options
	Cache    gouache.Cache
	Database gouache.Database
}

func New(c gouache.Cache, d gouache.Database, opts ...Option) gouache.Cache {
	return &cache{Options: newOptions(opts...), Cache: c, Database: d}
}

func (cache *cache) Get(ctx context.Context, key string) (any, error) {
	val, err := cache.Cache.Get(ctx, key)
	if errors.Is(err, gouache.ErrCacheMiss) {
		val, err := cache.Database.Select(ctx, key)
		if err != nil {
			return nil, err
		}
		return val, cache.Cache.Set(ctx, key, val)
	}
	return val, err
}

func (cache *cache) Set(ctx context.Context, key string, val any) error {
	if err := cache.Cache.Delete(ctx, key); err != nil {
		return err
	}
	if err := cache.Database.Upsert(ctx, key, val); err != nil {
		return err
	}
	return cache.Options.Gopher(func() {
		time.Sleep(cache.Options.DelayDuration)
		ctx := context.WithoutCancel(ctx)
		ctx, cancel := context.WithTimeout(ctx, cache.Options.DeleteTimeout)
		defer cancel()
		if err := cache.Cache.Delete(ctx, key); err != nil {
			cache.Options.ErrorHandler(err)
		}
	})
}

func (cache *cache) Delete(ctx context.Context, key string) error {
	if err := cache.Cache.Delete(ctx, key); err != nil {
		return err
	}
	if err := cache.Database.Delete(ctx, key); err != nil {
		return err
	}
	return cache.Options.Gopher(func() {
		time.Sleep(cache.Options.DelayDuration)
		ctx := context.WithoutCancel(ctx)
		ctx, cancel := context.WithTimeout(ctx, cache.Options.DeleteTimeout)
		defer cancel()
		if err := cache.Cache.Delete(ctx, key); err != nil {
			cache.Options.ErrorHandler(err)
		}
	})
}
