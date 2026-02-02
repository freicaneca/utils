package cache

import (
	"context"
	"testing"
	"time"
	"utils/logging"
	"utils/utils/testutils"
)

func TestCache(t *testing.T) {

	l := logging.New()

	ctx := context.Background()

	t.Run("remove matching", func(t *testing.T) {

		// will save two caches with matching keys
		// and another that doesnt match.
		// will then remove the matching ones:
		// only one cache must remain.

		h := NewCacheRAM()

		e := &Entry{
			Key:  "baba",
			Data: "bobo",
			CreationDateTime: time.Date(
				2025, 9, 10, 10, 0, 0, 0, time.UTC,
			),
			ExpirationDateTime: time.Date(
				2049, 9, 10, 10, 0, 0, 0, time.UTC,
			),
		}

		err := h.SaveCache(l, ctx, e)
		testutils.AssertError(t, err, nil)

		e = &Entry{
			Key:  "kababaka",
			Data: "bobo",
			CreationDateTime: time.Date(
				2025, 9, 10, 10, 0, 0, 0, time.UTC,
			),
			ExpirationDateTime: time.Date(
				2049, 9, 10, 10, 0, 0, 0, time.UTC,
			),
		}

		err = h.SaveCache(l, ctx, e)
		testutils.AssertError(t, err, nil)

		e = &Entry{
			Key:  "kokoko",
			Data: "bobo",
			CreationDateTime: time.Date(
				2025, 9, 10, 10, 0, 0, 0, time.UTC,
			),
			ExpirationDateTime: time.Date(
				2049, 9, 10, 10, 0, 0, 0, time.UTC,
			),
		}
		err = h.SaveCache(l, ctx, e)
		testutils.AssertError(t, err, nil)

		h.RemoveMatchingCaches(l, ctx, "baba")

		_, err = h.GetCache(
			l, ctx, "baba", 443322*time.Hour,
		)
		testutils.AssertError(t, err, ErrNotFound)

		_, err = h.GetCache(
			l, ctx, "kababaka", 443322*time.Hour,
		)
		testutils.AssertError(t, err, ErrNotFound)

		got, err := h.GetCache(
			l, ctx, "kokoko", 443322*time.Hour,
		)
		testutils.AssertError(t, err, nil)

		testutils.AssertStruct(
			t, got, e,
		)

	})

	t.Run("save, get, old entry", func(t *testing.T) {

		h := NewCacheRAM()

		e := &Entry{
			Key:  "baba",
			Data: "bobo",
			CreationDateTime: time.Date(
				2025, 9, 10, 10, 0, 0, 0, time.UTC,
			),
			ExpirationDateTime: time.Now().Add(200 * time.Second),
		}

		err := h.SaveCache(l, ctx, e)
		testutils.AssertError(t, err, nil)

		// will sleep for a while so the cache gets old
		time.Sleep(2 * time.Second)

		// max age required: 1 second. cache is old
		got, err := h.GetCache(
			l, ctx, "baba", 1*time.Second,
		)
		testutils.AssertError(t, err, ErrOlderThanMaxAge)
		testutils.AssertBool(t, got == nil, true)

	})

	t.Run("save, get, expire, get, not found", func(t *testing.T) {

		h := NewCacheRAM()

		e := &Entry{
			Key:  "baba",
			Data: "bobo",
			CreationDateTime: time.Date(
				2025, 9, 10, 10, 0, 0, 0, time.UTC,
			),
			ExpirationDateTime: time.Now().Add(2 * time.Second),
		}

		err := h.SaveCache(l, ctx, e)
		testutils.AssertError(t, err, nil)

		got, err := h.GetCache(
			l, ctx, "baba", 443322*time.Hour,
		)
		testutils.AssertError(t, err, nil)

		got.ExpirationDateTime = e.ExpirationDateTime
		testutils.AssertStruct(
			t, got, e,
		)

		time.Sleep(3 * time.Second)

		// now, cache should have expired

		got, err = h.GetCache(
			l, ctx, "baba", 443322*time.Hour,
		)
		testutils.AssertError(t, err, ErrNotFound)
		testutils.AssertBool(t, got == nil, true)

	})

	t.Run("save, get", func(t *testing.T) {

		h := NewCacheRAM()

		e := &Entry{
			Key:  "baba",
			Data: "bobo",
			CreationDateTime: time.Date(
				2025, 9, 10, 10, 0, 0, 0, time.UTC,
			),
			ExpirationDateTime: time.Date(
				2049, 9, 10, 10, 0, 0, 0, time.UTC,
			),
		}

		err := h.SaveCache(l, ctx, e)
		testutils.AssertError(t, err, nil)

		got, err := h.GetCache(
			l, ctx, "baba", 443322*time.Hour,
		)
		testutils.AssertError(t, err, nil)

		testutils.AssertStruct(
			t, got, e,
		)

	})

	t.Run("get from empty cache", func(t *testing.T) {

		h := NewCacheRAM()

		got, err := h.GetCache(
			l, ctx, "baba", 1*time.Second)
		testutils.AssertError(t, err, ErrNotFound)
		testutils.AssertBool(t, got == nil, true)

	})

}
