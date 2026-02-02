package cryptoutils

import (
	"crypto/hmac"
	"crypto/sha256"
)

func HMAC(
	payload []byte,
	secret []byte,
) []byte {

	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)

	return mac.Sum(nil)
}
