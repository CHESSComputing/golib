package auth

import (
	"errors"
	"log"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Response struct {
	Status string `json:"status"`
	Uid    int    `json:"uid,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

// Token represents access token structure
type Token struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (t *Token) Validate(clientId string) error {
	// validate our token
	var jwtKey = []byte(clientId)
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(t.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return errors.New("invalid signature")
			log.Fatal(err)
		}
		return err
	}
	if !tkn.Valid {
		return errors.New("invalid token")
	}
	return nil
}
