package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
	services "github.com/CHESSComputing/golib/services"
	"github.com/gin-gonic/gin"
)

// _token is used across all authorized APIs
var _token *Token

// gin cookies
// https://gin-gonic.com/docs/examples/cookie/
// more advanced use-case:
// https://stackoverflow.com/questions/66289603/use-existing-session-cookie-in-gin-router
func TokenMiddleware(clientId string, verbose int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if user request has valid token
		tokenStr := RequestToken(c.Request)
		token := &Token{AccessToken: tokenStr}
		if err := token.Validate(clientId); err != nil {
			msg := fmt.Sprintf("TokenMiddleware: invalid token %s, error %v", tokenStr, err)
			log.Println("ERROR:", msg)
			rec := services.Response("authz", http.StatusUnauthorized, services.TokenError, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, rec)
			return
		}
		if verbose > 0 {
			log.Println("INFO: token is validated")
		}
		c.Next()
	}
}

// ScopeTokenMiddleware provides token validation with specific scope
func ScopeTokenMiddleware(scope, clientId string, verbose int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if user request has valid token
		tokenStr := RequestToken(c.Request)
		token := &Token{AccessToken: tokenStr}
		if err := token.Validate(clientId); err != nil {
			msg := fmt.Sprintf("ScopeTokenMiddleware: invalid token %s, error %v", tokenStr, err)
			log.Println("ERROR:", msg)
			rec := services.Response("authz", http.StatusUnauthorized, services.TokenError, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, rec)
			return
		}
		if verbose > 0 {
			log.Println("INFO: write token is validated")
		}
		// check if token has proper write scope
		claims, err := TokenClaims(tokenStr, srvConfig.Config.Authz.ClientID)
		if err != nil {
			msg := fmt.Sprintf("ScopeTokenMiddleware: token '%s' error '%s'", tokenStr, err)
			log.Println("ERROR:", msg)
			log.Println("token", tokenStr)
			rec := services.Response("authz", http.StatusUnauthorized, services.ScopeError, errors.New(msg))
			c.AbortWithStatusJSON(http.StatusUnauthorized, rec)
			return
		}
		tscope := claims.CustomClaims.Scope
		if !strings.Contains(tscope, scope) {
			msg := fmt.Sprintf("ScopeTokenMiddleware: token scope '%s' does not match with scope '%s'", tscope, scope)
			log.Println("ERROR:", msg)
			log.Println("token", tokenStr)
			rec := services.Response("authz", http.StatusUnauthorized, services.ScopeError, errors.New(msg))
			c.AbortWithStatusJSON(http.StatusUnauthorized, rec)
			return
		}
		c.Next()
	}
}

// RequestToken gets token from http request
func RequestToken(r *http.Request) string {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return tokenStr
	}
	arr := strings.Split(tokenStr, " ")
	token := arr[len(arr)-1]
	return token
}
