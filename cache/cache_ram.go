package cache

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	"utils/logging"
)

type cacheRam struct {
	// key: entry key
	data map[string]*Entry
	lock *sync.Mutex
}

func NewCacheRAM() *cacheRam {
	return &cacheRam{
		data: map[string]*Entry{},
		lock: &sync.Mutex{},
	}
}

func (h *cacheRam) SaveCache(
	log *logging.Logger,
	ctx context.Context,
	entry *Entry,
) error {

	l := log.New()

	if entry == nil {
		return fmt.Errorf("%w: null entry",
			ErrBadRequest)
	}

	if !entry.ExpirationDateTime.IsZero() &&
		entry.ExpirationDateTime.Before(time.Now()) {
		l.Warn("expiration dt before now. do nothin")
		return nil
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	h.data[entry.Key] = entry

	l.Debug("registered cache with key %v",
		entry.Key)

	l.Debug("will expire in %v",
		entry.ExpirationDateTime.Format(time.RFC3339Nano))

	if !entry.ExpirationDateTime.IsZero() {
		go time.AfterFunc(
			time.Until(entry.ExpirationDateTime),
			func() {
				h.lock.Lock()
				defer h.lock.Unlock()

				delete(h.data, entry.Key)

				l.Debug("entry %v expired", entry.Key)
			})
	}

	return nil
}

func (h *cacheRam) GetCache(
	log *logging.Logger,
	ctx context.Context,
	key string,
	maxAge time.Duration,
) (
	*Entry,
	error,
) {

	l := log.New()

	h.lock.Lock()
	defer h.lock.Unlock()

	entry, ok := h.data[key]
	if !ok {
		return nil, ErrNotFound
	}

	if maxAge > 0 &&
		time.Since(entry.CreationDateTime) >
			maxAge {
		l.Warn("old cache")
		return nil, ErrOlderThanMaxAge
	}

	return entry, nil
}

func (h *cacheRam) RemoveCache(
	log *logging.Logger,
	ctx context.Context,
	key string,
) {

	l := log.New()

	h.lock.Lock()
	defer h.lock.Unlock()

	_, ok := h.data[key]
	if !ok {
		l.Debug("cache %v not found to be delete",
			key)
		return
	}

	delete(h.data, key)

	l.Debug("cache %v deleted", key)

}

func (h *cacheRam) RemoveMatchingCaches(
	log *logging.Logger,
	ctx context.Context,
	matchKey string,
) {

	l := log.New()

	h.lock.Lock()
	defer h.lock.Unlock()
	for k := range h.data {
		if strings.Contains(k, matchKey) {
			delete(h.data, k)
			l.Debug("cache %v deleted", k)
		}
	}

}
