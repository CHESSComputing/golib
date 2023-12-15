package server

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
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
