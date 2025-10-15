// Package gocache provides an implementation of the gouache.Cache interface
// using patrickmn/go-cache as the underlying storage mechanism.
//
// This package enables in-memory caching with expiration capabilities by leveraging
// go-cache's thread-safe operations and automatic expiration handling.
package gocache

import (
	"context"
	"time"

	"github.com/go-leo/gouache"
	"github.com/patrickmn/go-cache"
)

// Ensure that Cache implements the gouache.Cache interface at compile time.
var _ gouache.Cache = (*Cache)(nil)

// Cache is an implementation of gouache.Cache using go-cache as the storage backend.
// It provides methods for storing, retrieving, and deleting cached values with
// support for configurable time-to-live (TTL) settings.
type Cache struct {
	// Cache is the underlying go-cache instance used for storage.
	Cache *cache.Cache

	// TTL is an optional function to determine the time-to-live duration for a cache entry.
	// If not provided, the default expiration behavior of go-cache is used.
	TTL func(ctx context.Context, key string, val any) (time.Duration, error)
}

// Get retrieves a value from the cache by its key.
// It returns gouache.ErrNil if the key does not exist or has expired.
//
// Parameters:
//   - ctx: Context for the operation
//   - key: The key to retrieve the value for
//
// Returns:
//   - The cached value or nil if not found
//   - An error if the operation fails, or gouache.ErrNil if key doesn't exist
func (store *Cache) Get(ctx context.Context, key string) (any, error) {
	// Attempt to get the value from the go-cache
	val, ok := store.Cache.Get(key)
	
	// Handle case where entry is not found or has expired
	if !ok {
		return nil, gouache.ErrNil
	}
	
	// Return the found value
	return val, nil
}

// Set stores a value in the cache with the given key.
// The TTL function can be used to determine a custom expiration time for the entry.
//
// Parameters:
//   - ctx: Context for the operation
//   - key: The key to store the value under
//   - val: The value to store
//
// Returns:
//   - An error if the TTL function fails, nil otherwise
func (store *Cache) Set(ctx context.Context, key string, val any) error {
	// If a TTL function is provided, use it to determine expiration
	if store.TTL != nil {
		ttl, err := store.TTL(ctx, key, val)
		if err != nil {
			return err
		}
		// Set the value with the custom TTL
		store.Cache.Set(key, val, ttl)
		return nil
	}
	
	// Use default expiration behavior
	store.Cache.Set(key, val, cache.DefaultExpiration)
	return nil
}

// Delete removes a value from the cache by its key.
//
// Parameters:
//   - ctx: Context for the operation
//   - key: The key of the value to delete
//
// Returns:
//   - Always returns nil as go-cache.Delete doesn't return errors
func (store *Cache) Delete(ctx context.Context, key string) error {
	// Delegate deletion to the underlying go-cache instance
	store.Cache.Delete(key)
	return nil
}