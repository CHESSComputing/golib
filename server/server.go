package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	authz "github.com/CHESSComputing/golib/authz"
	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-gonic/gin"
)

var _routes gin.RoutesInfo

// Route represents routes structure
type Route struct {
	Method     string
	Path       string
	Handler    gin.HandlerFunc
	Authorized bool
}

// Router provids server router, it takes two maps:
// one for non-authorized routes and anotehr for authorized ones
func Router(routes []Route, fsys fs.FS, static, base string, verbose int) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// GET routes
	r.GET("/apis", APIsHandler)

	// loop over routes and creates necessary router structure
	var authGroup bool
	for _, route := range routes {
		if route.Authorized {
			authGroup = true
			continue
		}
		if route.Method == "GET" {
			r.GET(route.Path, route.Handler)
		} else if route.Method == "POST" {
			r.POST(route.Path, route.Handler)
		} else if route.Method == "PUT" {
			r.PUT(route.Path, route.Handler)
		} else if route.Method == "DELETE" {
			r.DELETE(route.Path, route.Handler)
		}
	}

	// all authorized routes
	if authGroup {
		authorized := r.Group("/")
		authorized.Use(authz.TokenMiddleware(srvConfig.Config.Authz.ClientID, verbose))
		{
			for _, route := range routes {
				if !route.Authorized {
					continue
				}
				if route.Method == "GET" {
					authorized.GET(route.Path, route.Handler)
				} else if route.Method == "POST" {
					authorized.POST(route.Path, route.Handler)
				} else if route.Method == "PUT" {
					authorized.PUT(route.Path, route.Handler)
				} else if route.Method == "DELETE" {
					authorized.DELETE(route.Path, route.Handler)
				}
			}
		}
	}

	// static files
	if fsys != nil {
		if entries, err := os.ReadDir(static); err == nil {
			for _, e := range entries {
				dir := e.Name()
				filesFS, err := fs.Sub(fsys, filepath.Join(static, dir))
				if err != nil {
					panic(err)
				}
				m := fmt.Sprintf("%s/%s", base, dir)
				r.StaticFS(m, http.FS(filesFS))
			}
		}
	}
	_routes = r.Routes()

	return r
}
