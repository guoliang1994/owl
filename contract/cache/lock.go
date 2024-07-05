package cache

import (
	"context"
)

type Lock interface {
	// Get Attempt to acquire the lock.
	Get(ctx context.Context, callback func()) error

	// Block Attempt to acquire the lock for the given number of seconds.
	Block(ctx context.Context, seconds int, callback func()) error

	// Release the lock.
	Release() bool

	// Owner Returns the current owner of the lock.
	Owner() string

	// ForceRelease Releases this lock in disregard of ownership.
	ForceRelease()
}
