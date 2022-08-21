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

type Server struct {
	Handler404 gin.HandlerFunc
}

var APP_PORT string = config.AppSetting.AppPort
var WSS_PREFIX string = config.WssSetting.Prefix
var APP_COOKIE *config.Cookie = config.CookieSetting

var defaultServer = &Server{
	Handler404: func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404", gin.H{
			"content": "Page not found",
		})
	},
}

func Setup(starterServer *Server) *gin.Engine {

	copier.CopyWithOption(&defaultServer, &starterServer, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// new engine
	engine := gin.New()

	store := persistence.NewInMemoryStore(time.Second)

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

	// apply middleware
	if APP_COOKIE.Enable {
		cookieStore := cookie.NewStore([]byte(APP_COOKIE.Secret))
		cookieStore.Options(sessions.Options{MaxAge: 60 * 1}) // expire in one minute
		engine.Use(sessions.Sessions(APP_COOKIE.Name, cookieStore))
	}

	engine.Use(middleware.Recovery(recoveryHandler)) // error nice handle
	engine.Use(static.Serve("/static", static.LocalFile("./static", false)))
	engine.Use(favicon.New("./static/favicon.ico")) // set favicon middleware
	engine.Use(middleware.CORS(middleware.CORSOptions{}))
	engine.Use(middleware.AccessFormat())

	// info router
	engine.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "info", "message": "Ok"})
	})

	engine.GET("/info_cache", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": "hello world, " + fmt.Sprint(time.Now().Unix()), "message": "Ok"})
	}))

	// index page
	engine.GET("", controller.IndexRouter)
	engine.GET("/index", controller.IndexRouter)
	engine.GET("/home", controller.IndexRouter)

	// about page
	engine.GET("/about", controller.AboutRouter)

	// labs
	engine.GET("/labs-bigint", actions.BigRouter)
	engine.GET("/labs-ytld", actions.YtIdRouter)
	engine.GET("/labs-pem", actions.PemRouter)

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
	c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": errors.Wrap(err, 3).Error()})
}
