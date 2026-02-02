package core

import (
	"testing"
	"utils/logging"
	"utils/utils/stringutils"
	"utils/utils/testutils"
)

func TestCore(t *testing.T) {

	l := logging.New()

	t.Run("remove", func(t *testing.T) {

		h := New(
			stringutils.RandomString(4), 1, 5,
			//func(arg *Req) bool { return false },
		)

		go h.Run(l, func(arg *Req) bool {
			return false
		})

		h.PushBack("1", "3")
		h.PushBack("2", "3")

		testutils.AssertInt(
			t, h.elements.Len(), 2,
		)

		h.Remove(l, "1")

		testutils.AssertInt(
			t, h.elements.Len(), 1,
		)

	})

	t.Run("push back", func(t *testing.T) {

		h := New(
			stringutils.RandomString(4), 1, 5,
			//func(arg *Req) bool { return false },
		)

		go h.Run(l, func(arg *Req) bool {
			return false
		})

		h.PushBack("1", "3")
		h.PushBack("2", "3")

		testutils.AssertInt(
			t, h.elements.Len(), 2,
		)

	})

}
