package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"utils/logging"
)

func EncryptAES(
	log *logging.Logger,
	secret []byte,
	key []byte,
) (
	string,
	error,
) {

	l := log.New()

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	for i := 0; i < gcm.NonceSize(); i++ {
		nonce[i] = 'x'
	}

	encrypted := gcm.Seal(nil, nonce, secret, nil)

	out := base64.StdEncoding.EncodeToString(encrypted)

	l.Debug("encrypted secret *** with key %v: %v",
		string(key), out)

	return out, nil
}

func DecryptAES(
	log *logging.Logger,
	encryptedB64 string,
	key []byte,
) (
	[]byte,
	error,
) {

	l := log.New()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	for i := 0; i < gcm.NonceSize(); i++ {
		nonce[i] = 'x'
	}

	encrypted, err := base64.StdEncoding.DecodeString(encryptedB64)
	if err != nil {
		return nil, fmt.Errorf("decoding b64: %v", err)
	}

	secret, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm open: %v", err)
	}

	l.Debug("decrypted %v with key %v", encryptedB64, string(key))

	return secret, nil
}
