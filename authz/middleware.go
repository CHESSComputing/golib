package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
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

// ScopeTokenMiddleware provides token validation with specific scope
func ScopeTokenMiddleware(scope, clientId string, verbose int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if user request has valid token
		tokenStr := RequestToken(c.Request)
		token := &Token{AccessToken: tokenStr}
		if err := token.Validate(clientId); err != nil {
			msg := fmt.Sprintf("TokenMiddleware: invalid token %s, error %v", tokenStr, err)
			log.Println("ERROR:", msg)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"status": "fail", "error": err.Error()})
			return
		}
		if verbose > 0 {
			log.Println("INFO: write token is validated")
		}
		// check if token has proper write scope
		if token.Scope != scope {
			msg := fmt.Sprintf("ScopeTokenMiddleware: token scope '%s' does not match with scope '%s'", token.Scope, scope)
			log.Println("ERROR:", msg)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"status": "fail", "error": errors.New(msg)})
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

// helper function to refresh global token used in authorized APIs
func refreshToken() error {
	// check and obtain token
	var err error
	if _token == nil {
		if token, err := getToken(); err == nil {
			_token = &token
		} else {
			return err
		}
	} else {
		err = _token.Validate(srvConfig.Config.Authz.ClientID)
	}
	return err
}

// helper function to obtain JWT token from Authz service
func getToken() (Token, error) {
	var token Token
	rurl := fmt.Sprintf(
		"%s/oauth/token?client_id=%s&client_secret=%s&grant_type=client_credentials&scope=read",
		srvConfig.Config.Services.AuthzURL,
		srvConfig.Config.Authz.ClientID,
		srvConfig.Config.Authz.ClientSecret)
	resp, err := http.Get(rurl)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return token, err
	}
	reqToken := token.AccessToken
	if srvConfig.Config.Authz.Verbose > 0 {
		log.Printf("INFO: obtain token %+v", token)
	}

	// validate our token
	var jwtKey = []byte(srvConfig.Config.Authz.ClientID)
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return token, errors.New("invalid token signature")
		}
		return token, err
	}
	if !tkn.Valid {
		log.Println("WARNING: token invalid")
		return token, errors.New("invalid token validity")
	}
	return token, nil
}

// gin cookies
// https://gin-gonic.com/docs/examples/cookie/
// more advanced use-case:
// https://stackoverflow.com/questions/66289603/use-existing-session-cookie-in-gin-router
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if our user key is set
		if user, err := c.Cookie("user"); err == nil {
			if srvConfig.Config.Authz.Verbose > 0 {
				log.Println(c.Request.Method, c.Request.URL.Path, user)
			}
			c.Set("user", user)
			if err := refreshToken(); err != nil {
				c.SetCookie("user", "", -1, "/", "localhost", false, true)
				c.Set("user", "")
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte("Unauthorized access"))

				//                 content := errorTmpl(c, "unable to get valid token", err)
				//                 log.Fatal(content)

				//                 c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte(content))
				//                 c.Redirect(http.StatusFound, "/")
				return
			}
			return
		}

		if user, ok := c.Get("user"); !ok {
			if srvConfig.Config.Authz.Verbose > 0 {
				log.Println(c.Request.Method, c.Request.URL.Path)
			}
			c.Redirect(http.StatusFound, "/login")
		} else {
			if srvConfig.Config.Authz.Verbose > 0 {
				log.Println(c.Request.Method, c.Request.URL.Path, user)
			}
		}
		if err := refreshToken(); err != nil {
			//             content := errorTmpl(c, "unable to get valid token", err)
			c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte("Unauthorized access"))
			return
		}
		c.Next()
	}
}
