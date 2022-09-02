package starter

import (
	"fmt"
	"html/template"
	"net/http"
	_ "strings"
	"time"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/gredis"
	"github.com/chunhui2001/go-starter/wss"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/static"
	"github.com/go-errors/errors"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/middleware"
	_ "github.com/chunhui2001/go-starter/mycache"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/chunhui2001/go-starter/controller"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/jinzhu/copier"
	"github.com/thinkerou/favicon"
)

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
}

var store *persistence.InMemoryStore = persistence.NewInMemoryStore(time.Second)
var APP_PORT string = config.AppSetting.AppPort
var WSS_PREFIX string = config.WssSetting.Prefix
var APP_COOKIE *config.Cookie = config.CookieSetting
var WEB_PAGE_CONF *config.WebPageConf = config.WebPageSettings

func rateLimitKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

var defaultServer = &Server{
	Store: store,
	HandlerInfo: func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "info", "message": "Ok"})
	},
	HandlerIndexPage: controller.IndexRouter,
	Handler404: func(c *gin.Context) {
		if WEB_PAGE_CONF.Enable {
			c.HTML(http.StatusNotFound, "404", gin.H{
				"requestUrl": c.Request.URL.Path,
				"content":    "Page not found",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusNotFound,
				"message": "Page not found",
			})
		}
	},
	Handler500: func(c *gin.Context, err interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": errors.Wrap(err, 3).Error()})
	},
	CustomeRoutes: []Route{
		{Method: http.MethodGet, Path: "/info_cache", Handlers: []gin.HandlerFunc{
			cache.CachePage(store, time.Minute, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"code": 200, "data": "hello world, " + fmt.Sprint(time.Now().Unix()), "message": "Ok"})
			}),
		}},
		{Method: http.MethodGet, Path: "/about", Handlers: []gin.HandlerFunc{controller.AboutRouter}},
	},
}

func Setup(starterServer *Server) *gin.Engine {

	copier.CopyWithOption(&defaultServer, &starterServer, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	if config.AppSetting.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// new engine
	engine := gin.New()

	if WEB_PAGE_CONF.Enable {
		// init html template
		engine.HTMLRender = ginview.New(goview.Config{
			Root:      WEB_PAGE_CONF.Root,
			Extension: WEB_PAGE_CONF.Extension,
			Master:    WEB_PAGE_CONF.Master,
			//Partials:  []string{"partials/ad"},
			Funcs: template.FuncMap{
				"sub": func(a, b int) int {
					return a - b
				},
				// more funcs
			},
			DisableCache: true,
		})
	}

	// cookie session
	if APP_COOKIE.Enable {
		cookieStore := cookie.NewStore([]byte(APP_COOKIE.Secret))
		cookieStore.Options(sessions.Options{MaxAge: 60 * 1}) // expire in one minute
		engine.Use(sessions.Sessions(APP_COOKIE.Name, cookieStore))
	}

	// RateLimit: This makes it so each ip can only make 5 requests per second
	ratelimitMiddleWare := ratelimit.RateLimiter(
		// ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		// 	Rate:  time.Second,
		// 	Limit: 5,
		// }),
		ratelimit.RedisStore(&ratelimit.RedisOptions{
			RedisClient: gredis.Client(),
			Rate:        time.Second,
			Limit:       5,
		}),
		&ratelimit.Options{
			ErrorHandler: errorHandler,
			KeyFunc:      rateLimitKeyFunc,
		})

	// apply middleware
	engine.Use(middleware.Recovery(recoveryHandler)) // error nice handle
	engine.Use(static.Serve("/static", static.LocalFile("./static", false)))
	engine.Use(favicon.New("./static/favicon.ico")) // set favicon middleware
	engine.Use(middleware.CORS(middleware.CORSOptions{}))
	engine.Use(middleware.AccessLog())

	if WEB_PAGE_CONF.Enable {

		engine.GET("/info", ratelimitMiddleWare, defaultServer.HandlerInfo) // info router
		engine.GET("", ratelimitMiddleWare, defaultServer.HandlerIndexPage) // index page
		engine.GET("/index", ratelimitMiddleWare, defaultServer.HandlerIndexPage)
		engine.GET("/home", ratelimitMiddleWare, defaultServer.HandlerIndexPage)

		if WEB_PAGE_CONF.LoginUrl != "" {
			engine.GET(WEB_PAGE_CONF.LoginUrl, controller.LoginHandler)
		}

		if WEB_PAGE_CONF.SignUpUrl != "" {
			engine.GET(WEB_PAGE_CONF.SignUpUrl, controller.SignUpHandler)
		}

	}

	// customer routes
	for _, ro := range defaultServer.CustomeRoutes {
		engine.Handle(ro.Method, ro.Path, ro.Handlers...)
	}

	if WSS_PREFIX != "" {
		engine.GET(WSS_PREFIX, wss.WebsocketUpgrade)
	}

	engine.NoRoute(func(c *gin.Context) {
		if c.Request.RequestURI == "/favicon.ico" {
			c.Next()
		} else {
			defaultServer.Handler404(c)
		}
	})

	config.Log.Info("Listening and serving HTTP on " + APP_PORT + ", websocket=" + config.WssSetting.Wss())

	return engine

}

func recoveryHandler(c *gin.Context, err interface{}) {
	defaultServer.Handler500(c, err)
}
