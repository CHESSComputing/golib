package server

import (
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
