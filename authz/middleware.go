package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
		tokenStr := getToken(c.Request)
		token := &Token{AccessToken: tokenStr}
		if err := token.Validate(clientId); err != nil {
			msg := fmt.Sprintf("invalid token %s, error %v", tokenStr, err)
			log.Println("WARNING:", msg)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"status": "fail", "error": err.Error()})
			return
		}
		if verbose > 0 {
			log.Println("INFO: token is validated")
		}
		c.Next()
	}
}

// helper function to get token from http request
func getToken(r *http.Request) string {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return tokenStr
	}
	arr := strings.Split(tokenStr, " ")
	token := arr[len(arr)-1]
	return token
}
