package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	rsxid "github.com/rs/xid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var github_config *oauth2.Config

type githubUser struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTML_URL          string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Name              string `json:"name"`
	Company           string `json:"company"`
	Blog              string `json:"blog"`
	Location          string `json:"location"`
	Email             string `json:"email"`
	Hireable          bool   `json:"hireable"`
	Bio               string `json:"bio"`
	TwitterUserName   string `json:"twitter_username"`
	PublicRepos       int    `json:"public_repos"`
	PublicGists       int    `json:"public_gits"`
	Followers         int    `json:"followers"`
	Following         int    `json:"following"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// helper function to get github oauth config and state
func getGithubOauthURL() (*oauth2.Config, string) {
	for _, arec := range srvConfig.Config.Frontend.OAuth {
		if arec.Provider == "github" {
			github_config = &oauth2.Config{
				ClientID:     arec.ClientID,
				ClientSecret: arec.ClientSecret,
				RedirectURL:  arec.RedirectURL,
				Scopes:       []string{"user", "repo"},
				Endpoint:     github.Endpoint,
			}
		}
	}
	state := rsxid.New().String()
	return github_config, state
}

// GithubOauthLogin provides gin handler for github oauth login
func GithubOauthLogin(ctx *gin.Context, verbose int) {
	config, state := getGithubOauthURL()
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

// GithubCallBack provides gin handler for github callback to given endpoint
func GithubCallBack(ctx *gin.Context, endpoint string, verbose int) {
	session := sessions.Default(ctx)
	state := session.Get("state")
	if state != ctx.Query("state") {
		msg := "GithubCallBack state error"
		_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New(msg))
		return
	}

	code := ctx.Query("code")
	token, err := github_config.Exchange(ctx, code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	client := github_config.Client(context.TODO(), token)
	userInfo, err := client.Get("https://api.github.com/user")
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

	var user githubUser
	err = json.Unmarshal(info, &user)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if verbose > 0 {
		log.Printf("### github json %+v", user)
	}
	// set necessary cookie for our web server
	//     expiration := time.Now().Add(24 * time.Hour)
	//     cookie := http.Cookie{Name: "user", Value: user.Login, Expires: expiration}
	//     http.SetCookie(ctx.Writer, &cookie)
	ctx.Set("user", user.Login)
	ctx.SetCookie("user", user.Login, 7200, "/", domain(), false, true)
	ctx.Redirect(http.StatusSeeOther, endpoint)
}

// helper function to get host domain
func domain() string {
	domain := "localhost"
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: unable to get hostname, error:", err)
	}
	if !strings.Contains(hostname, ".") {
		hostname = "localhost"
	} else {
		arr := strings.Split(hostname, ".")
		domain = strings.Join(arr[len(arr)-2:], ".")
	}
	return domain
}
