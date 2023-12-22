package server

import (
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// Time0 represents initial time when we started the server
var Time0 time.Time

// init function
func init() {
	Time0 = time.Now()
}

//
// GET handlers
//

// CaptchaHandler provides access to captcha server
func CaptchaHandler() gin.HandlerFunc {
	hdlr := captcha.Server(captcha.StdWidth, captcha.StdHeight)
	return func(c *gin.Context) {
		hdlr.ServeHTTP(c.Writer, c.Request)
	}
}

// GinRoute represents git route info
type GinRoute struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// ApisHandler provides JSON output for server routes
func ApisHandler(c *gin.Context, routes gin.RoutesInfo) {
	var ginRoutes []GinRoute
	for _, r := range routes {
		route := GinRoute{Method: r.Method, Path: r.Path}
		ginRoutes = append(ginRoutes, route)
	}
	c.JSON(http.StatusOK, ginRoutes)
}
