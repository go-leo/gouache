// Package gouache provides cache and database interfaces for building cache systems.
//
// This package defines the core interfaces that cache implementations should follow,
// enabling consistent integration with various storage backends.
package gouache

import (
	"context"
	"errors"
)

// ErrCacheMiss represents a cache miss error, returned when a requested key
// does not exist in the cache.
var ErrCacheMiss = errors.New("gouache: key not found")

// Cache defines the basic operations for a cache implementation.
type Cache interface {
	// Get retrieves a value from the cache by its key.
	// It returns ErrCacheMiss if the key does not exist.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - key: The key to retrieve the value for
	//
	// Returns:
	//   - The cached value or nil if not found
	//   - An error if the operation fails, or ErrCacheMiss if key doesn't exist
	Get(ctx context.Context, key string) (any, error)

	// Set stores a value in the cache under the specified key.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - key: The key under which the value will be stored
	//   - val: The value to store
	//
	// Returns:
	//   - An error if the operation fails
	Set(ctx context.Context, key string, val any) error

	// Delete removes a value from the cache by its key.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - key: The key of the value to delete
	//
	// Returns:
	//   - An error if the operation fails
	Delete(ctx context.Context, key string) error
}

// Database defines the basic operations for a database implementation.
type Database interface {
	// Select retrieves a record from the database by its key.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - key: The key to query the record for
	//
	// Returns:
	//   - The queried record or nil if not found
	//   - An error if the operation fails
	Select(ctx context.Context, key string) (any, error)

	// Upsert inserts or updates a record in the database.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - key: The key of the record to upsert
	//   - val: The value to store
	//
	// Returns:
	//   - An error if the operation fails
	Upsert(ctx context.Context, key string, val any) error

	// Delete removes a record from the database by its key.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - key: The key of the record to delete
	//
	// Returns:
	//   - An error if the operation fails
	Delete(ctx context.Context, key string) error
}
