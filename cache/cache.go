package cache

import (
	"context"
	"errors"
	"time"
	"utils/logging"
)

type Handler interface {
	SaveCache(
		log *logging.Logger,
		ctx context.Context,
		entry *Entry,
	) error

	// gets cache not older than maxAge.
	// if maxAge = 0, no age checking.
	GetCache(
		log *logging.Logger,
		ctx context.Context,
		key string,
		maxAge time.Duration,
	) (
		*Entry,
		error,
	)

	RemoveCache(
		log *logging.Logger,
		ctx context.Context,
		key string,
	)

	// Removes all caches whose key contains matchKey.
	RemoveMatchingCaches(
		log *logging.Logger,
		ctx context.Context,
		matchKey string,
	)
}

type Entry struct {
	// usually, unique string for a given data.
	Key              string
	Data             string
	CreationDateTime time.Time
	// if zero, no expiration
	ExpirationDateTime time.Time
}

var (
	ErrInternal        = errors.New("internal err")
	ErrOlderThanMaxAge = errors.New("older than max age")
	ErrBadRequest      = errors.New("bad request")
	ErrNotFound        = errors.New("not found")
)
