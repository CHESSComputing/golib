package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	authz "github.com/CHESSComputing/golib/authz"
	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logging "github.com/vkuznet/http-logging"
)

var _routes gin.RoutesInfo

// Route represents routes structure
type Route struct {
	Method     string
	Path       string
	Handler    gin.HandlerFunc
	Authorized bool
}

// InitServer provides server initialization
func InitServer(webServer srvConfig.WebServer) {
	// setup log options
	rotateLogs(webServer.LogFile)

	// setup log options
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if webServer.LogLongFile {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}

	// setup gin options
	if webServer.GinOptions.DisableConsoleColor {
		// Disable Console Color
		gin.DisableConsoleColor()
	}

	// set gin mode of operation
	if webServer.GinOptions.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if webServer.GinOptions.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if webServer.GinOptions.Mode == "test" {
		gin.SetMode(gin.TestMode)
	}

	// use production flag to overwrite gin mode
	if webServer.GinOptions.Production {
		gin.SetMode(gin.ReleaseMode)
	}
	log.Printf("webServer configuration:\n%s", webServer.String())
}

// Router provids server router, it takes two maps:
// one for non-authorized routes and anotehr for authorized ones
func Router(routes []Route, fsys fs.FS, static string, webServer srvConfig.WebServer) *gin.Engine {
	base := webServer.Base
	verbose := webServer.Verbose

	InitServer(webServer)

	// setup gin router
	//     r := gin.Default()
	r := gin.New()

	// GET routes
	r.GET("/apis", ApisHandler)

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

// helper function to rotate logs
func rotateLogs(srvLogName string) {
	if srvLogName != "" {
		log.SetOutput(new(logging.LogWriter))
		rl, err := rotatelogs.New(logName(srvLogName))
		if err == nil {
			rotlogs := logging.RotateLogWriter{RotateLogs: rl}
			log.SetOutput(rotlogs)
		}
	}
}

// logName returns proper log name based on Config LogFile and either
// hostname or pod name (used in k8s environment).
func logName(srvLogName string) string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("unable to get hostname", err)
	}
	if os.Getenv("MY_POD_NAME") != "" {
		hostname = os.Getenv("MY_POD_NAME")
	}
	logName := srvLogName + "_%Y%m%d"
	if hostname != "" {
		logName = fmt.Sprintf("%s_%s", srvLogName, hostname) + "_%Y%m%d"
	}
	return logName
}
