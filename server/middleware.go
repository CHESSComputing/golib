package server

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

// TotalGetRequests counts total number of GET requests received by the server
var TotalGetRequests uint64

// TotalPostRequests counts total number of POST requests received by the server
var TotalPostRequests uint64

// TotalPutRequests counts total number of PUT requests received by the server
var TotalPutRequests uint64

// TotalDeleteRequests counts total number of DELETE requests received by the server
var TotalDeleteRequests uint64

// CounterMiddleware counts GET/POST/PUT/DELETE requests
func CounterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "GET" {
			atomic.AddUint64(&TotalGetRequests, 1)
		} else if method == "POST" {
			atomic.AddUint64(&TotalPostRequests, 1)
		} else if method == "PUT" {
			atomic.AddUint64(&TotalPutRequests, 1)
		} else if method == "DELETE" {
			atomic.AddUint64(&TotalDeleteRequests, 1)
		}
		c.Next()
	}
}

// LimiterMiddleware provides limiter middleware pointer
var LimiterMiddleware gin.HandlerFunc

// initialize Limiter middleware pointer
func initLimiter(period, header string) {
	log.Printf("limiter rate='%s'", period)
	// create rate limiter with 5 req/second
	rate, err := limiter.NewRateFromFormatted(period)
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	if header != "" {
		instance = limiter.New(
			store,
			rate,
			limiter.WithClientIPHeader(header))
	}
	LimiterMiddleware = mgin.NewMiddleware(instance)
}

// helper function to get hash of the string, provided by https://github.com/amalfra/etag
func getHash(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

// Generates an Etag for given string, provided by https://github.com/amalfra/etag
func Etag(str string, weak bool) string {
	tag := fmt.Sprintf("\"%d-%s\"", len(str), getHash(str))
	if weak {
		tag = "W/" + tag
	}
	return tag
}

// HeaderMiddleware represents header middleware
func HeaderMiddleware(webServer srvConfig.WebServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		goVersion := runtime.Version()
		tstamp := time.Now().Format("2006-02-01")
		server := fmt.Sprintf("foxden (%s %s)", goVersion, tstamp)
		w.Header().Add("Server", server)

		// settng Etag and its expiration
		if r.Method == "GET" && webServer.Etag != "" && webServer.CacheControl != "" {
			etag := Etag(webServer.Etag, false)
			w.Header().Set("Etag", etag)
			w.Header().Set("Cache-Control", webServer.CacheControl) // 5 minutes
			if match := r.Header.Get("If-None-Match"); match != "" {
				if strings.Contains(match, etag) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}
		c.Next()
	}
}

// LoggerMiddleware is custom logger for gin server
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process the request
		c.Next()

		// Calculate the duration
		duration := time.Since(start).Seconds()

		// Get the status code that is sent to the client
		statusCode := c.Writer.Status()

		// Get the client IP address
		clientIP := c.ClientIP()

		// Log the request details
		r := c.Request

		// Save the current log flags
		originalFlags := log.Flags()

		// Set custom log flags
		log.SetFlags(log.Ldate | log.Ltime)

		// yield log message about HTTP request
		log.Printf("%s %d %s %s [client: %s] [%v sec]",
			r.Proto,
			statusCode,
			r.Method,
			c.Request.URL.Path,
			clientIP,
			duration)

		// Restore the original log flags
		log.SetFlags(originalFlags)
	}
}
