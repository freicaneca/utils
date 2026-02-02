package cache

import (
	"context"
	"testing"
	"time"
	"utils/logging"
	"utils/utils/testutils"
)

func TestUtil(t *testing.T) {

	l := logging.New()

	ctx := context.Background()

	t.Run("save, get", func(t *testing.T) {

		h := NewCacheRAM()

		type person struct {
			Name string
			Age  int
		}

		data := person{
			Name: "baba",
			Age:  15,
		}

		err := SaveFromModel(
			l, ctx, h, "baba", time.Date(
				2049, 9, 10, 10, 0, 0, 0, time.UTC,
			), data,
		)
		testutils.AssertError(t, err, nil)

		out := person{}

		err = GetToModel(
			l, ctx, h, "baba", 443322*time.Hour, &out,
		)
		testutils.AssertError(t, err, nil)

		testutils.AssertStruct(
			t, out, data,
		)

	})

}
