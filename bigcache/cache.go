// Package bigcache provides an implementation of the gouache.Cache interface
// using allegro/bigcache as the underlying storage mechanism.
//
// This package enables high-performance caching capabilities by leveraging
// BigCache's efficient memory management and concurrent access patterns.
package bigcache

import (
	"context"
	"errors"

	"github.com/allegro/bigcache/v3"
	"github.com/go-leo/gouache"
)

// Ensure that Cache implements the gouache.Cache interface at compile time.
var _ gouache.Cache = (*Cache)(nil)

// Cache is an implementation of gouache.Cache using BigCache as the storage backend.
// It provides methods for storing, retrieving, and deleting cached values with
// support for custom serialization and deserialization functions.
type Cache struct {
	// BigCache is the underlying BigCache instance used for storage.
	BigCache *bigcache.BigCache

	// Marshal is an optional function to serialize objects into bytes.
	// If not provided, default type conversions are used for basic types.
	Marshal func(key string, obj any) ([]byte, error)

	// Unmarshal is an optional function to deserialize bytes into objects.
	// If not provided, raw bytes are returned.
	Unmarshal func(key string, data []byte) (any, error)
}

// Get retrieves a value from the cache by its key.
// It returns gouache.ErrNil if the key does not exist.
//
// Parameters:
//   - ctx: Context for the operation
//   - key: The key to retrieve the value for
//
// Returns:
//   - The cached value or nil if not found
//   - An error if the operation fails, or gouache.ErrNil if key doesn't exist
func (store *Cache) Get(ctx context.Context, key string) (any, error) {
	// Attempt to get the value from BigCache
	data, err := store.BigCache.Get(key)

	// Handle case where entry is not found
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		return nil, gouache.ErrNil
	}

	// Return other errors as-is
	if err != nil {
		return nil, err
	}

	// If no unmarshal function is defined, return raw data
	if store.Unmarshal == nil {
		return data, nil
	}

	// Use custom unmarshal function to decode the data
	obj, err := store.Unmarshal(key, data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Set stores a value in the cache with the given key.
// It supports various basic types when no Marshal function is provided.
//
// Parameters:
//   - ctx: Context for the operation
//   - key: The key to store the value under
//   - val: The value to store
//
// Returns:
//   - An error if the operation fails
func (store *Cache) Set(ctx context.Context, key string, val any) error {
	// If no custom marshal function is provided, use built-in type handling
	if store.Marshal == nil {
		switch val := val.(type) {
		case []byte:
			return store.BigCache.Set(key, val)
		default:
			return errors.New("bigcache: val is not []byte")
		}
	}

	// Use custom marshal function to encode the value
	data, err := store.Marshal(key, val)
	if err != nil {
		return err
	}

	// Store the encoded data in BigCache
	return store.BigCache.Set(key, data)
}

// Delete removes a value from the cache by its key.
//
// Parameters:
//   - ctx: Context for the operation
//   - key: The key of the value to delete
//
// Returns:
//   - An error if the operation fails
func (store *Cache) Delete(ctx context.Context, key string) error {
	// Delegate deletion to the underlying BigCache instance
	return store.BigCache.Delete(key)
}