package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"
	"utils/logging"
)

func SaveFromModel(
	log *logging.Logger,
	ctx context.Context,
	cache Handler,
	key string,
	expDT time.Time,
	data any,
) error {

	l := log.New()

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return fmt.Errorf("encoding data: %v: %w",
			err, ErrInternal)
	}

	nowDT := time.Now()
	err = cache.SaveCache(
		l, ctx, &Entry{
			Key:                key,
			Data:               buf.String(),
			CreationDateTime:   nowDT,
			ExpirationDateTime: expDT,
		},
	)
	if err != nil {
		return fmt.Errorf("saving to cache: %w",
			err)
	}

	l.Debug("saved to cache")

	return nil
}

// out must be an address of the desired type.
func GetToModel(
	log *logging.Logger,
	ctx context.Context,
	cache Handler,
	key string,
	maxAge time.Duration,
	out any,
) error {

	l := log.New()

	got, err := cache.GetCache(
		l, ctx, key, maxAge,
	)
	if err != nil {
		return fmt.Errorf("getting cache: %w", err)
	}

	l.Debug("fetching from cache")

	dec := gob.NewDecoder(
		bytes.NewReader([]byte(got.Data)),
	)

	err = dec.Decode(out)
	if err != nil {
		return fmt.Errorf("decoding from cache: %v: %w",
			err, ErrInternal)
	}

	return nil
}
