package auth

import (
	"encoding/hex"
	"errors"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// type Response struct {
//     Status string `json:"status"`
//     Uid    int    `json:"uid,omitempty"`
//     Error  string `json:"error,omitempty"`
// }

// CustomClaims defines application specific claims
type CustomClaims struct {
	User        string   `json:"user"`
	Scope       string   `json:"scope"`
	Kind        string   `json:"kind"`
	Roles       []string `json:"roles"`
	Application string   `json:"application"`
}

// Claims defines our JWT claims
type Claims struct {
	jwt.RegisteredClaims
	CustomClaims CustomClaims `json:"custom_claims"`
}

// Token represents access token structure
type Token struct {
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// Validate performs token validation
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

// TokenClaims returns token claims
func TokenClaims(accessToken, clientId string) *Claims {
	var jwtKey = []byte(clientId)
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("ERROR", err)
		}
	}
	if !tkn.Valid {
		err := errors.New("invalid token")
		log.Println("ERROR", err)
	}
	return claims
}

// JWTAccessToken generates JWT access token with custom claims
// https://blog.canopas.com/jwt-in-golang-how-to-implement-token-based-authentication-298c89a26ffd
func JWTAccessToken(secretKey string, expiresAt int64, customClaims CustomClaims) (string, error) {
	var sub, aud string
	if uuid, err := uuid.NewRandom(); err == nil {
		sub = hex.EncodeToString(uuid[:])
	}
	if uuid, err := uuid.NewRandom(); err == nil {
		aud = hex.EncodeToString(uuid[:])
	}
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "CHESS Authz server",
			// the `sub` (Subject) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
			Subject: sub,

			// the `aud` (Audience) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
			Audience: jwt.ClaimStrings{aud},

			// the `exp` (Expiration Time) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresAt) * time.Second)),

			// the `nbf` (Not Before) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
			//             NotBefore *NumericDate `json:"nbf,omitempty"`

			// the `iat` (Issued At) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
			IssuedAt: jwt.NewNumericDate(time.Now()),

			// the `jti` (JWT ID) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
			//             ID string `json:"jti,omitempty"`
		},
		CustomClaims: customClaims,
	}

	// generate a string using claims and HS256 algorithm
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// sign the generated key using secretKey
	// SignedString declared as interface{} but should accept []byte
	// see https://github.com/dgrijalva/jwt-go/issues/65
	token, err := tokenString.SignedString([]byte(secretKey))

	return token, err
}
