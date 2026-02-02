package cryptoutils

import (
	"testing"
	"utils/logging"
	"utils/utils/testutils"
)

func TestCrypt(t *testing.T) {

	l := logging.New()

	t.Run("encrypt, decrypt", func(t *testing.T) {

		secret := "bububababababubu"
		key := "kakakokokakakoko"

		encr, err := EncryptAES(
			l, []byte(secret), []byte(key),
		)
		testutils.AssertError(t, err, nil)

		got, err := DecryptAES(
			l, encr, []byte(key),
		)
		testutils.AssertError(t, err, nil)
		testutils.AssertString(t, string(got), secret)

	})

}
