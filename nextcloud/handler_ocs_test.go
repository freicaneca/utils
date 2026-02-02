package nextcloud

import (
	"context"
	"testing"
	"utils/logging"
	"utils/utils/testutils"
)

func TestHandler(t *testing.T) {

	l := logging.New()

	ctx := context.Background()

	url := "https://cloud.dev.10x.umnicorn.com"
	admin := "nc_admin"
	pw := "Pp7pa#7w$)(&hr6}+JNCh#5q"

	h, err := NewHandlerOCS(
		url, admin, pw,
	)
	testutils.AssertError(t, err, nil)

	t.Run("register, remove", func(t *testing.T) {

		err := h.RegisterUser(
			l,
			ctx,
			"patati@patata.com",
			"4jngnejgej085j",
		)
		testutils.AssertError(t, err, nil)

		err = h.RemoveUser(
			l, ctx, "patati@patata.com",
		)
		testutils.AssertError(t, err, nil)

	})

}
