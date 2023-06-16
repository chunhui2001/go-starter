package starter

import (
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/chunhui2001/go-starter/controller"
	"github.com/chunhui2001/go-starter/core"
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/gproxy"
	"github.com/chunhui2001/go-starter/core/gredis"
	"github.com/chunhui2001/go-starter/core/gwss"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/chunhui2001/go-starter/graph"
	"github.com/chunhui2001/go-starter/graph/generated"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/go-errors/errors"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/core/middleware"
	_ "github.com/chunhui2001/go-starter/mycache"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/jinzhu/copier"
)

type ReverseProxy struct {
	From    string   `yaml:"from"`
	To      string   `yaml:"to"`
	Remotes []string `yaml:"remotes"`
}

type Route struct {
	Method   string
	Path     string
	Handlers []gin.HandlerFunc
}

type Server struct {
	Store            *persistence.InMemoryStore
	HandlerInfo      gin.HandlerFunc
	HandlerIndexPage gin.HandlerFunc
	Handler404       gin.HandlerFunc
	Handler500       func(c *gin.Context, err interface{})
	CustomeRoutes    []Route
	R                *gin.Engine
	GraphqlHandler   gin.HandlerFunc
}

func rateLimitKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {

	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", graphServerConf.ServerURi)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var defaultServer = &Server{
	Store: store,
	HandlerInfo: func(c *gin.Context) {
		c.JSON(http.StatusOK, (&R{Data: fmt.Sprintf(`Yeah, your server %s is running.`, APP_SETTINGS.AppName)}).Success())
	},
	HandlerIndexPage: controller.IndexRouter,
	Handler404: func(c *gin.Context) {
		if WEB_PAGE_CONF.Enable {
			c.HTML(http.StatusNotFound, "404", gin.H{
				"requestUrl": c.Request.URL.Path,
				"content":    "Page not found",
			})
		} else {
			c.JSON(http.StatusOK, (&R{}).Fail(404))
		}
	},
	Handler500: func(c *gin.Context, err interface{}) {
		c.JSON(http.StatusInternalServerError, (&R{Error: errors.Wrap(err, 3)}).Fail(500))
	},
	CustomeRoutes: []Route{
		{Method: http.MethodGet, Path: "/info_cache", Handlers: []gin.HandlerFunc{
			cache.CachePage(store, time.Minute, func(c *gin.Context) {
				c.JSON(http.StatusOK, (&R{Data: "hello world, " + fmt.Sprint(time.Now().Unix())}).Success())
			}),
		}},
	},
	GraphqlHandler: graphqlHandler(),
}

var (
	APP_SETTINGS      *config.AppConf            = config.AppSetting
	APP_PORT          string                     = config.AppSetting.AppPort
	WSS_Conf          *config.Wss                = config.WssSetting
	APP_COOKIE        *config.Cookie             = config.CookieSetting
	WEB_PAGE_CONF     *config.WebPageConf        = config.WebPageSettings
	store             *persistence.InMemoryStore = persistence.NewInMemoryStore(time.Second)
	Redis_Conf        *gredis.GRedis             = config.RedisConf
	logger                                       = config.Log
	graphServerConf                              = config.GraphServerSetting
	reverseProxyArray []ReverseProxy
)

func (s *Server) Bootstrap(hooks ...func(*gin.Engine)) *Server {

	if Redis_Conf.Mode != gredis.Disabled {
		for _, channel := range strings.Split(Redis_Conf.SubChannels, ",") {
			if Redis_Conf.PrintMessage {
				gredis.Sub(channel, func(channel string, payload string) {
					config.Log.Info("RedisMessage-收到了消息: channel=" + channel + ", payload=" + payload)
				})
			}
		}
	}

	if err := copier.CopyWithOption(&defaultServer, &s, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		panic(err)
	}

	s.R = Setup()

	for _, h := range hooks {
		h(s.R)
	}

	return s

}

func (s *Server) Running() {

	srv := &http.Server{
		Addr:        APP_SETTINGS.AppPort,
		Handler:     s.R,
		IdleTimeout: 5 * time.Second,
	}

	// srv.SetKeepAlivesEnabled(true) // 默认是true

	utils.AddShutDownHook(config.Log, func() {
		config.Log.Info("shutting down server")
		// clean up
		if err := srv.Shutdown(context.Background()); err != nil {
			config.Log.Info("shutting down server-err")
		} else {
			config.Log.Info("shutting down server-done")
		}
	})

	l, err := net.Listen("tcp", APP_SETTINGS.AppPort)

	if err != nil {
		config.Log.Info("Application Run Failed: ErrorMessage=" + err.Error())
		os.Exit(1)
		return
	}

	go func() {
		config.Log.Info("Congratulations! Your server startup successfully, Listening and serving HTTP on " + APP_SETTINGS.AppPort)
		config.Log.Info(srv.Serve(l))
	}()

	utils.WaitShutDown()

}

func (s *Server) RunningTLS() {

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:         APP_SETTINGS.AppPort,
		Handler:      s.R,
		IdleTimeout:  5 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// srv.SetKeepAlivesEnabled(true) // 默认是true

	utils.AddShutDownHook(config.Log, func() {
		config.Log.Info("shutting down server")
		// clean up
		if err := srv.Shutdown(context.Background()); err != nil {
			config.Log.Info("shutting down server-err")
		} else {
			config.Log.Info("shutting down server-done")
		}
		config.Log.Info("shutting down server-done")
	})

	l, err := net.Listen("tcp", APP_SETTINGS.AppPort)

	if err != nil {
		config.Log.Info("Application Run Failed: ErrorMessage=" + err.Error())
		os.Exit(1)
		return
	}

	go func() {
		config.Log.Infof(`Congratulations! Your server startup successfully, Listening and serving HTTP on %s`, APP_SETTINGS.AppPort)
		config.Log.Info(srv.ServeTLS(l, "server.crt", "server.key"))
	}()

	utils.WaitShutDown()
}

func ginFuncMap() template.FuncMap {

	funcMaps := template.FuncMap{
		"string": func(b any) string {
			return utils.ToString(b)
		},
		"plainstring": func(b any) string {
			return fmt.Sprintf("%.0f", b)
		},
		"timestring": func(b uint32) string {
			return time.Unix(int64(b), 0).Format("2006-01-02T15:04:05Z07:00")
		},
		"GIN_MAPS_ENV": func(b string) string {
			return config.GetEnv(b, "")
		},
		// more funcs
	}

	return funcMaps

}

func Setup() *gin.Engine {

	if APP_SETTINGS.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// new engine
	engine := gin.New()

	if WEB_PAGE_CONF.Enable {
		// https://curatedgo.com/r/goview-is-a-foolingoview/index.html
		// https://noknow.info/it/go/how_to_use_if_in_html_template?lang=ja
		// init html template
		engine.HTMLRender = ginview.New(goview.Config{
			Root:         filepath.Join(config.AppRoot(), WEB_PAGE_CONF.Root),
			Extension:    WEB_PAGE_CONF.Extension,
			Master:       WEB_PAGE_CONF.Master,
			Partials:     []string{"partials/ad"},
			Funcs:        ginFuncMap(),
			DisableCache: true,
		})
	}

	var rateLimitStore ratelimit.Store = ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 5,
	})

	if !Redis_Conf.Disabled() {
		rateLimitStore = ratelimit.RedisStore(&ratelimit.RedisOptions{
			RedisClient: gredis.Client(),
			Rate:        time.Second,
			Limit:       5,
		})
	}

	// RateLimit: This makes it so each ip can only make 5 requests per second
	ratelimitMiddleWare := ratelimit.RateLimiter(
		rateLimitStore,
		&ratelimit.Options{
			ErrorHandler: errorHandler,
			KeyFunc:      rateLimitKeyFunc,
		})

	// cookie session
	if APP_COOKIE.Enable {
		store := cookie.NewStore([]byte(APP_COOKIE.Secret))
		store.Options(sessions.Options{MaxAge: APP_COOKIE.MaxAge})
		engine.Use(sessions.Sessions(APP_COOKIE.Name, store))
	}

	if err := config.ReadConfig("Reverse-Proxy", &reverseProxyArray); err != nil {
		logger.Errorf("Reverse-Proxy-Configuration-Error: Key=%s, ErrorMessage=%s", "Reverse-Proxy", err.Error())
	}

	for _, val := range reverseProxyArray {
		gproxy.Any(engine, val.From, val.To, val.Remotes...)
	}

	// apply middlewares
	engine.Use(middleware.Urlwriter())               // urlwriter
	engine.Use(middleware.Recovery(recoveryHandler)) // error nice handle
	engine.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".pdf", ".mp4", ".ico"})))

	if ok, _ := utils.FileExists(filepath.Join(config.AppRoot(), "static")); ok {
		engine.Use(static.Serve("/static", static.LocalFile(filepath.Join(config.AppRoot(), "./static"), false)))
	}

	engine.Use(core.Favicon(filepath.Join(config.AppRoot(), "./static/favicon.ico"))) // set favicon middleware
	engine.Use(middleware.CORS(middleware.CORSOptions{}))
	engine.Use(middleware.AccessLog("/favicon.ico", "/static", "/info"))

	// default info
	engine.GET("/info", defaultServer.HandlerInfo) // info router

	if WEB_PAGE_CONF.Enable {

		// builtin pages
		AppendRouter(http.MethodGet, []string{"/", "/index", "home"}, ratelimitMiddleWare, defaultServer.HandlerIndexPage)
		AppendRouter(http.MethodGet, []string{"/about"}, ratelimitMiddleWare, controller.AboutRouter)

		if WEB_PAGE_CONF.LoginUrl != "" {
			AppendRouter(http.MethodGet, []string{WEB_PAGE_CONF.LoginUrl}, controller.LoginHandler)
			AppendRouter(http.MethodPost, []string{WEB_PAGE_CONF.LoginUrl}, controller.PostLoginHandler)
			AppendRouter(http.MethodGet, []string{"/logout"}, controller.LogoutHandler)
		}

		if WEB_PAGE_CONF.SignUpUrl != "" {
			AppendRouter(http.MethodGet, []string{WEB_PAGE_CONF.SignUpUrl}, controller.SignUpHandler)
		}

	}

	// REGISTER ROUTES
	for _, ro := range defaultServer.CustomeRoutes {
		logger.Infof("REGISTER-A-ROUTER: Method=%s, Path=%s, Handlers=%s", ro.Method, ro.Path, JoinHandlersString(ro.Handlers))
		engine.Handle(ro.Method, ro.Path, ro.Handlers...)
	}

	if WSS_Conf.Enable {
		engine.GET(WSS_Conf.Prefix, gwss.WebsocketUpgrade)
	}

	engine.NoRoute(func(c *gin.Context) {
		if c.Request.RequestURI == "/favicon.ico" {
			c.Next()
		} else {
			defaultServer.Handler404(c)
		}
	})

	if WSS_Conf.Enable {
		logger.Info("Startup a websocket server running on " + config.WssSetting.Wss())
	}

	if graphServerConf.Enable {
		engine.POST(graphServerConf.ServerURi, defaultServer.GraphqlHandler)
		engine.GET(graphServerConf.PlayGroundURi, playgroundHandler())
	}

	return engine

}

func AppendRouter(method string, paths []string, handlers ...gin.HandlerFunc) {
	for _, path := range paths {
		defaultServer.CustomeRoutes = append(defaultServer.CustomeRoutes, Route{Method: method, Path: path, Handlers: handlers})
	}
}

func JoinHandlersString(handlers []gin.HandlerFunc) string {
	var handlersString []string
	for _, h := range handlers {
		handlersString = append(handlersString, utils.GetFunctionName(h))
	}
	return strings.Join(handlersString, ", ")
}

func recoveryHandler(c *gin.Context, err interface{}) {
	defaultServer.Handler500(c, err)
}
