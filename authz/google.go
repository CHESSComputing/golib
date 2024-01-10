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
	"golang.org/x/oauth2/google"
)

var google_config *oauth2.Config

type googleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

// helper function to get facebook oauth config and state
func getGoogleOauthURL() (*oauth2.Config, string) {
	for _, arec := range srvConfig.Config.Frontend.OAuth {
		if arec.Provider == "google" {
			google_config = &oauth2.Config{
				ClientID:     arec.ClientID,
				ClientSecret: arec.ClientSecret,
				RedirectURL:  arec.RedirectURL,
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint: google.Endpoint,
			}
		}
	}
	state := rsxid.New().String()
	return google_config, state
}

// GoogleOauthLogin provides gin handler for google oauth login
func GoogleOauthLogin(ctx *gin.Context, verbose int) {
	config, state := getGoogleOauthURL()
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

// GoogleCallBack provides gin handler for google callback to given endpoint
func GoogleCallBack(ctx *gin.Context, endpoint string, verbose int) {
	session := sessions.Default(ctx)
	state := session.Get("state")
	if state != ctx.Query("state") {
		msg := "GoogleCallBack state error"
		_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New(msg))
		return
	}

	code := ctx.Query("code")
	token, err := google_config.Exchange(ctx, code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	client := google_config.Client(context.TODO(), token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer userInfo.Body.Close()

	info, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user googleUser
	err = json.Unmarshal(info, &user)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if verbose > 0 {
		log.Printf("### google json %+v", user)
	}
	// set necessary cookie for our web server
	ctx.Set("user", user.Name)
	ctx.SetCookie("user", user.Name, 7200, "/", domain(), false, true)
	ctx.Redirect(http.StatusSeeOther, endpoint)
}
