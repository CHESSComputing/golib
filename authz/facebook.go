package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	rsxid "github.com/rs/xid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var facebook_config *oauth2.Config

type facebookUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// helper function to get facebook oauth config and state
func getFacebookOauthURL() (*oauth2.Config, string) {
	for _, arec := range srvConfig.Config.Frontend.OAuth {
		if arec.Provider == "facebook" {
			facebook_config = &oauth2.Config{
				ClientID:     arec.ClientID,
				ClientSecret: arec.ClientSecret,
				RedirectURL:  arec.RedirectURL,
				Scopes: []string{
					"email",
					"public_profile",
				},
				Endpoint: facebook.Endpoint,
			}
		}
	}
	state := rsxid.New().String()
	return facebook_config, state
}

// FacebookOauthLogin provides gin handler for facebook oauth login
func FacebookOauthLogin(ctx *gin.Context, verbose int) {
	config, state := getFacebookOauthURL()
	if verbose > 0 {
		log.Printf("GithubOauthLogin config %+v", config)
	}
	redirectURL := config.AuthCodeURL(state)
	session := sessions.Default(ctx)
	session.Set("state", state)
	err := session.Save()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusSeeOther, redirectURL)
}

// FacebookCallBack provides gin handler for facebook callback to given endpoint
func FacebookCallBack(ctx *gin.Context, endpoint string, verbose int) {
	if error_reason := ctx.Query("error_reason"); error_reason != "" {
		_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New(error_reason))
		return
	}

	session := sessions.Default(ctx)
	state := session.Get("state")
	if state != ctx.Query("state") {
		msg := "FacebookCallBack state error"
		_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New(msg))
		return
	}

	code := ctx.Query("code")
	token, err := facebook_config.Exchange(ctx, code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	client := facebook_config.Client(context.TODO(), token)

	userInfo, err := client.Get("https://graph.facebook.com/v8.0/me?fields=id,name,email")
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	info, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user facebookUser
	err = json.Unmarshal(info, &user)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if verbose > 0 {
		log.Printf("### facebook json %+v", user)
	}
	// set necessary cookie for our web server
	ctx.Set("user", user.Name)
	ctx.SetCookie("user", user.Name, 7200, "/", domain(), false, true)
	ctx.Redirect(http.StatusSeeOther, endpoint)
}
