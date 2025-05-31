package auth

import (
	"testing"
)

// TestToken
func TestToken(t *testing.T) {
	scope := "read"
	customClaims := CustomClaims{Scope: scope}
	secretKey := "lksjdlfkjsd"
	expiresAt := int64(100)
	tokenStr, err := JWTAccessToken(secretKey, expiresAt, customClaims)
	if err != nil {
		t.Error(err)
	}
	token := Token{AccessToken: tokenStr, Expires: expiresAt, Scope: scope}
	if err := token.Validate(secretKey); err != nil {
		t.Errorf("invalid token")
	}
	claims, err := TokenClaims(tokenStr, secretKey)
	if claims.CustomClaims.Scope != scope {
		t.Errorf("invalid scope")
	}
	if err != nil {
		t.Error(err)
	}
}
