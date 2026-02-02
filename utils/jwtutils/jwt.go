package jwtutils

import (
	"net/http"
	"strings"
	"utils/logging"

	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateJWT(
	log *logging.Logger,
	rawToken string,
	rsaPublicKey []byte,
) error {

	l := log.New()

	// converting public key to a proper format for jwt library
	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKey)
	if err != nil {
		l.Error("error parsing rsa public key: %v", err)
		return err
	}

	// now parsing rawToken, creating a function that authenticates it
	// based on parsedKey
	_, err = jwt.Parse(rawToken, func(tk *jwt.Token) (
		any, error) {
		_, ok := tk.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, jwt.ErrNotRSAPublicKey
		}
		return parsedKey, nil
	})

	if err != nil {
		l.Error("error parsing token: %v", err)
		return err
	}

	return nil
}

func GetClaimFromRequest(
	r *http.Request,
	claim string,
) string {

	token := r.Header.Get("Authorization")

	fields := strings.Split(token, " ")

	if len(fields) != 2 {
		return ""
	}

	token = fields[1]

	claims, err := GetUnverifiedJWTClaims(
		logging.New(),
		token,
	)
	if err != nil {
		return ""
	}

	userID, ok := claims["user_id"]
	if !ok {
		return ""
	}

	return userID.(string)
}

func GetUnverifiedJWTClaims(
	log *logging.Logger,
	rawToken string,
) (
	jwt.MapClaims,
	error,
) {

	l := log.New()

	p := jwt.NewParser()

	token, _, err := p.ParseUnverified(rawToken, jwt.MapClaims{})
	if err != nil {
		l.Error("error parsing unverified raw token: %v", err)
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
