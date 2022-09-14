package gproxy

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/gin-gonic/gin"
)

var (
	logger = config.Log
)

func Proxy(prefix string, remotes string, c *gin.Context) {

	rand.Seed(time.Now().UnixNano())

	upstreams := strings.Split(remotes, ",")
	upstreamSize := len(upstreams)
	currentRemote := upstreams[rand.Intn((upstreamSize-1)-0+1)+0]

	upstream, err := url.Parse(currentRemote)

	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(upstream)

	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = upstream.Host
		req.URL.Scheme = upstream.Scheme
		req.URL.Host = upstream.Host
		req.URL.Path = prefix + "" + c.Param("proxyPath")
		logger.Infof(`Reverse-Proxy: Upstream=%s, ProxyPath=%s`, currentRemote, prefix+""+c.Param("proxyPath"))
	}

	proxy.ServeHTTP(c.Writer, c.Request)

}
