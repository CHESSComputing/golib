package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	authz "github.com/CHESSComputing/golib/authz"
	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logging "github.com/vkuznet/http-logging"
)

var _routes gin.RoutesInfo
var metricsPrefix string

// StartTime represents initial time when we started the server
var StartTime time.Time

// Route represents routes structure
type Route struct {
	Method     string
	Path       string
	Scope      string
	Authorized bool
	Handler    gin.HandlerFunc
}

// StartServer starts HTTP(s) server
func StartServer(r *gin.Engine, webServer srvConfig.WebServer) {
	sport := fmt.Sprintf(":%d", webServer.Port)
	if webServer.ServerKey != "" {
		certFile := webServer.ServerCrt
		ckeyFile := webServer.ServerKey
		log.Println("Start HTTPs server on port", sport)
		r.RunTLS(sport, certFile, ckeyFile)
	} else {
		log.Println("Start HTTP server on port", sport)
		r.Run(sport)
	}
}

// InitServer provides server initialization
func InitServer(webServer srvConfig.WebServer) {
	StartTime = time.Now()
	// setup log options
	rotateLogs(webServer.LogFile)

	// setup log options
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if webServer.LogLongFile {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}

	// setup limiter
	if webServer.LimiterPeriod == "" {
		// default 100 request per second
		webServer.LimiterPeriod = "100-S"
	}
	initLimiter(webServer.LimiterPeriod, webServer.LimiterHeader)
	metricsPrefix = webServer.MetricsPrefix

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
	r := gin.New()

	// initialize cookie store (used by authz module and oauth)
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("server_session", store))

	// GET routes
	r.GET("/apis", ApisHandler)
	r.GET("/metrics", MetricsHandler)

	// loop over routes and creates necessary router structure
	var authGroup bool
	var readRoutes, writeRoutes []Route
	for _, route := range routes {
		if route.Authorized {
			authGroup = true
			if route.Scope == "write" {
				writeRoutes = append(writeRoutes, route)
			} else {
				readRoutes = append(readRoutes, route)
			}
			continue
		}
		log.Printf("method %s path %s auth %v scope '%s'", route.Method, route.Path, route.Authorized, route.Scope)
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
		authorizedRead := r.Group("/")
		authorizedRead.Use(authz.TokenMiddleware(srvConfig.Config.Authz.ClientID, verbose))
		{
			for _, route := range readRoutes {
				if !route.Authorized {
					continue
				}
				log.Printf("method %s path %s auth %v scope '%s'", route.Method, route.Path, route.Authorized, route.Scope)
				if route.Method == "GET" {
					authorizedRead.GET(route.Path, route.Handler)
				} else if route.Method == "POST" {
					authorizedRead.POST(route.Path, route.Handler)
				} else if route.Method == "PUT" {
					authorizedRead.PUT(route.Path, route.Handler)
				} else if route.Method == "DELETE" {
					authorizedRead.DELETE(route.Path, route.Handler)
				}
			}
		}
		authorizedWrite := r.Group("/")
		authorizedWrite.Use(authz.ScopeTokenMiddleware("write", srvConfig.Config.Authz.ClientID, verbose))
		{
			for _, route := range writeRoutes {
				if !route.Authorized {
					continue
				}
				log.Printf("method %s path %s auth %v scope '%s'", route.Method, route.Path, route.Authorized, route.Scope)
				if route.Method == "GET" {
					authorizedWrite.GET(route.Path, route.Handler)
				} else if route.Method == "POST" {
					authorizedWrite.POST(route.Path, route.Handler)
				} else if route.Method == "PUT" {
					authorizedWrite.PUT(route.Path, route.Handler)
				} else if route.Method == "DELETE" {
					authorizedWrite.DELETE(route.Path, route.Handler)
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

	// use common middlewares
	r.Use(CounterMiddleware())
	r.Use(LimiterMiddleware)
	r.Use(HeaderMiddleware(webServer))

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
		logName = fmt.Sprintf("%s_%s", hostname, srvLogName) + "_%Y%m%d"
	}
	return logName
}
