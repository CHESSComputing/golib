package server

import (
	"errors"
	"strings"

	authz "github.com/CHESSComputing/golib/authz"
	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-gonic/gin"
)

// GetAuthTokenUser returns user's token, user name and error
func GetAuthTokenUser(c *gin.Context) (string, string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", "", errors.New("no authorization header")
	}

	// Expect header format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", "", errors.New("no bearer token")
	}
	token := parts[1]
	claims, err := authz.TokenClaims(token, srvConfig.Config.Authz.ClientID)
	user := claims.CustomClaims.User

	return token, user, err
}
