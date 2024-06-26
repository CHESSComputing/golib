package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
func ApisHandler(c *gin.Context) {
	var ginRoutes []GinRoute
	for _, r := range _routes {
		route := GinRoute{Method: r.Method, Path: r.Path}
		ginRoutes = append(ginRoutes, route)
	}
	c.JSON(http.StatusOK, ginRoutes)
}

// MetricsHandler provides metrics JSON for monitoring purposes (Prometheus)
func MetricsHandler(c *gin.Context) {
	c.Writer.Write([]byte(promMetrics(metricsPrefix)))
}

// GinHandlerFunc converts given http.Handler to gin.HandlerFunc
func GinHandlerFunc(hdlr http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		hdlr(c.Writer, c.Request)
		c.Next()
	}
}

// QLKeysHandler provides list of keys used in QueryLanguage in given service
func QLKeysHandler(c *gin.Context) {
	var keys []string
	fname := fmt.Sprintf("%s/ql_keys.json", _staticDir)

	// read QL keys file
	file, err := os.Open(fname)
	if err != nil {
		log.Println("ERROR", err)
		c.JSON(http.StatusInternalServerError, keys)
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("ERROR", err)
		c.JSON(http.StatusInternalServerError, keys)
		return
	}

	// unmarshal our data into keys structure
	err = json.Unmarshal(data, &keys)
	if err != nil {
		log.Println("ERROR", err)
		c.JSON(http.StatusInternalServerError, keys)
		return
	}
	c.JSON(http.StatusOK, keys)
}
