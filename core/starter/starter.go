package starter

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/controller"
	"github.com/chunhui2001/go-starter/core"
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/config"
	"github.com/chunhui2001/go-starter/core/gredis"
	"github.com/chunhui2001/go-starter/core/grtask"
	"github.com/chunhui2001/go-starter/core/gwss"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
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

func rateLimitKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

var defaultServer = &Server{
	Store: store,
	HandlerInfo: func(c *gin.Context) {
		c.JSON(http.StatusOK, R{Data: "info"}.Success())
	},
	HandlerIndexPage: controller.IndexRouter,
	Handler404: func(c *gin.Context) {
		if WEB_PAGE_CONF.Enable {
			c.HTML(http.StatusNotFound, "404", gin.H{
				"requestUrl": c.Request.URL.Path,
				"content":    "Page not found",
			})
		} else {
			c.JSON(http.StatusOK, R{}.Fail(404))
		}
	},
	Handler500: func(c *gin.Context, err interface{}) {
		c.JSON(http.StatusInternalServerError, R{Error: errors.Wrap(err, 3)}.Fail(500))
	},
	CustomeRoutes: []Route{
		{Method: http.MethodGet, Path: "/info_cache", Handlers: []gin.HandlerFunc{
			cache.CachePage(store, time.Minute, func(c *gin.Context) {
				c.JSON(http.StatusOK, R{Data: "hello world, " + fmt.Sprint(time.Now().Unix())}.Success())
			}),
		}},
	},
}

var (
	APP_SETTINGS     *config.AppConf            = config.AppSetting
	APP_PORT         string                     = config.AppSetting.AppPort
	WSS_Conf         *config.Wss                = config.WssSetting
	APP_COOKIE       *config.Cookie             = config.CookieSetting
	WEB_PAGE_CONF    *config.WebPageConf        = config.WebPageSettings
	store            *persistence.InMemoryStore = persistence.NewInMemoryStore(time.Second)
	Redis_Conf       *gredis.GRedis             = config.RedisConf
	Simple_Task_Conf *config.SimpleGTask        = config.SimpleGTaskConf
	logger                                      = config.Log
)

func Bootstrap(starterServer *Server) *gin.Engine {

	if Redis_Conf.Mode != gredis.Disabled {
		for _, channel := range strings.Split(Redis_Conf.SubChannels, ",") {
			gredis.Sub(channel, func(channel string, payload string) {
				config.Log.Info("收到了消息1: channel=" + channel + ", payload=" + payload)
			})
		}
	}

	if starterServer != nil {
		copier.CopyWithOption(&defaultServer, &starterServer, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	}

	return Setup()

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

	// apply middleware
	engine.Use(middleware.Recovery(recoveryHandler)) // error nice handle

	if ok, _ := utils.FileExists(filepath.Join(utils.RootDir(), "static")); ok {
		engine.Use(static.Serve("/static", static.LocalFile("./static", false)))
	} else {
		logger.Warn("static folder not exists" + config.WssSetting.Wss())
	}

	engine.Use(core.Favicon("./static/favicon.ico")) // set favicon middleware
	engine.Use(middleware.CORS(middleware.CORSOptions{}))
	engine.Use(middleware.AccessLog())

	// default info
	engine.GET("/info", ratelimitMiddleWare, defaultServer.HandlerInfo) // info router

	if WEB_PAGE_CONF.Enable {

		// index page route
		AppendRouter(http.MethodGet, []string{"/", "/index", "home"}, ratelimitMiddleWare, defaultServer.HandlerIndexPage)
		// about simple page
		AppendRouter(http.MethodGet, []string{"/about"}, ratelimitMiddleWare, controller.AboutRouter)

		if WEB_PAGE_CONF.LoginUrl != "" {
			AppendRouter(http.MethodGet, []string{WEB_PAGE_CONF.LoginUrl}, controller.LoginHandler)
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

	logger.Info("Congratulations! Your server startup successfully, Listening and serving HTTP on " + APP_PORT)

	if WSS_Conf.Enable {
		logger.Info("Startup a websocket server running on " + config.WssSetting.Wss())
	}

	if Simple_Task_Conf.Enable {
		grtask.AddTask(APP_SETTINGS.AppName, Simple_Task_Conf.ID, Simple_Task_Conf.Name, Simple_Task_Conf.Expr, func(node string, taskId string) {
			for i := 0; i < 3; i++ {
				time.Sleep(1 * time.Second)
				config.Log.Infof("定时任务正在执行每秒1次,耗时3秒: num=%d, node=%s, taskId=%s", i+1, node, taskId)
			}
		})
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
