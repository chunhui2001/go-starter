package starter

import (
	"fmt"
	"html/template"
	"net/http"
	_ "strings"
	"time"

	"github.com/chunhui2001/go-starter/config"
	"github.com/chunhui2001/go-starter/wss"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/static"
	"github.com/go-errors/errors"

	"github.com/gin-gonic/gin"

	"github.com/chunhui2001/go-starter/logger"
	"github.com/chunhui2001/go-starter/middleware"
	_ "github.com/chunhui2001/go-starter/mycache"

	"github.com/chunhui2001/go-starter/actions"
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

var defaultServer = &Server{
	Store: store,
	HandlerInfo: func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "info", "message": "Ok"})
	},
	HandlerIndexPage: controller.IndexRouter,
	Handler404: func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404", gin.H{
			"content": "Page not found",
		})
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
		{Method: http.MethodGet, Path: "/labs-bigint", Handlers: []gin.HandlerFunc{actions.BigRouter}},
		{Method: http.MethodGet, Path: "/labs-ytld", Handlers: []gin.HandlerFunc{actions.YtIdRouter}},
		{Method: http.MethodGet, Path: "/labs-pem", Handlers: []gin.HandlerFunc{actions.PemRouter}},
	},
}

func Setup(starterServer *Server) *gin.Engine {

	copier.CopyWithOption(&defaultServer, &starterServer, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// new engine
	engine := gin.New()

	// init html template
	engine.HTMLRender = ginview.New(goview.Config{
		Root:      "views",
		Extension: ".html",
		Master:    "layouts/master",
		//Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
			// more funcs
		},
		DisableCache: true,
	})

	// cookie session
	if APP_COOKIE.Enable {
		cookieStore := cookie.NewStore([]byte(APP_COOKIE.Secret))
		cookieStore.Options(sessions.Options{MaxAge: 60 * 1}) // expire in one minute
		engine.Use(sessions.Sessions(APP_COOKIE.Name, cookieStore))
	}

	// apply middleware
	engine.Use(middleware.Recovery(recoveryHandler)) // error nice handle
	engine.Use(static.Serve("/static", static.LocalFile("./static", false)))
	engine.Use(favicon.New("./static/favicon.ico")) // set favicon middleware
	engine.Use(middleware.CORS(middleware.CORSOptions{}))
	engine.Use(middleware.AccessFormat())

	engine.GET("/info", defaultServer.HandlerInfo) // info router
	engine.GET("", defaultServer.HandlerIndexPage) // index page
	engine.GET("/index", defaultServer.HandlerIndexPage)
	engine.GET("/home", defaultServer.HandlerIndexPage)

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

	logger.Log.Info("Listening and serving HTTP on " + APP_PORT + ", websocket=" + config.WssSetting.Wss())

	return engine

}

func recoveryHandler(c *gin.Context, err interface{}) {
	defaultServer.Handler500(c, err)
}
