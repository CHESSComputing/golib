package server

import (
	"context"
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

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var _routes gin.RoutesInfo
var metricsPrefix string
var _staticDir string

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
	/*
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		if webServer.LogLongFile {
			log.SetFlags(log.LstdFlags | log.Llongfile)
		}
	*/
	// only use short/long file option of standard logger and do not use
	// log.LstdFlags (which is shortcut for log.Ldate | log.Ltime)
	// because we rely on our own logging module which provides timestamp
	log.SetFlags(log.Lshortfile)
	if webServer.LogLongFile {
		log.SetFlags(log.Llongfile)
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

	// initialize our web server
	InitServer(webServer)
	// remember server's static area (to be used in QLKeysHandler)
	_staticDir = static

	// setup gin router
	r := gin.New()

	// setup router logger middleware (it should comes before we setup routes)
	r.Use(LoggerMiddleware())

	// initialize cookie store (used by authz module and oauth)
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("server_session", store))

	// GET routes
	r.GET("/apis", ApisHandler)
	r.GET("/qlkeys", QLKeysHandler)
	r.GET("/metrics", MetricsHandler)

	// loop over routes and creates necessary router structure
	var authGroup bool
	var readRoutes, writeRoutes, deleteRoutes []Route
	for _, route := range routes {
		if route.Authorized {
			authGroup = true
			if route.Scope == "write" {
				writeRoutes = append(writeRoutes, route)
			} else if route.Scope == "delete" {
				deleteRoutes = append(deleteRoutes, route)
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
		//         authorizedRead.Use(authz.TokenMiddleware(srvConfig.Config.Authz.ClientID, verbose))
		authorizedRead.Use(authz.ScopeTokenMiddleware("read", srvConfig.Config.Authz.ClientID, verbose))
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
					authorizedRead.DELETE(route.Path, route.Handler)
				}
			}
		}
		authorizedDelete := r.Group("/")
		authorizedDelete.Use(authz.ScopeTokenMiddleware("delete", srvConfig.Config.Authz.ClientID, verbose))
		{
			for _, route := range deleteRoutes {
				if !route.Authorized {
					continue
				}
				log.Printf("method %s path %s auth %v scope '%s'", route.Method, route.Path, route.Authorized, route.Scope)
				if route.Method == "DELETE" {
					authorizedDelete.DELETE(route.Path, route.Handler)
				}
			}
		}
	}

	// static files
	if webServer.StaticDir != "" {
		log.Printf("Load static files from local area %s\n", webServer.StaticDir)
		if entries, err := os.ReadDir(webServer.StaticDir); err == nil {
			for _, e := range entries {
				dir := e.Name()
				m := fmt.Sprintf("%s/%s", base, dir)
				sdir := filepath.Join(webServer.StaticDir, dir)
				log.Printf("for end-point %s use static directory %s\n", m, sdir)
				r.StaticFS(m, http.Dir(sdir))
			}
		}
	} else if fsys != nil {
		if entries, err := fs.ReadDir(fsys, static); err == nil {
			for _, e := range entries {
				dir := e.Name()
				filesFS, err := fs.Sub(fsys, filepath.Join(static, dir))
				if err != nil {
					panic(err)
				}
				m := fmt.Sprintf("%s/%s", base, dir)
				if webServer.StaticDir != "" {
					sdir := filepath.Join(webServer.StaticDir, dir)
					log.Printf("for end-point %s use static directory %s\n", m, sdir)
					r.StaticFS(m, http.Dir(sdir))
				} else {
					log.Printf("for end-point %s use embeded fs\n", m)
					r.StaticFS(m, http.FS(filesFS))
				}
			}
		}
	} else {
		log.Println("WARNING: neither fsys or webServer.StaticDir is set, no css/js support")
	}

	_routes = r.Routes()

	// use common middlewares
	r.Use(CounterMiddleware())
	r.Use(LimiterMiddleware)
	r.Use(HeaderMiddleware(webServer))

	// open telemetry middlewares
	tp, err := InitTracer()
	if err == nil {
		defer func() {
			if tp != nil {
				_ = tp.Shutdown(context.Background())
			}
		}()
		r.Use(otelgin.Middleware("FOXDEN"))
		r.Use(TracingMiddleware())
	} else {
		log.Printf("WARNING: failed to initialize tracer: %v", err)
	}

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
