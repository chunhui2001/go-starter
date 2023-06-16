package gproxy

import (
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/config"
	"github.com/gin-gonic/gin"
)

var (
	logger                  = config.Log
	defaultTimeOut      int = 150 // * time.Second
	maxIdleConns        int = 100
	idleConnTimeout     int = 90
	maxIdleConnsPerHost int = 100
	maxConnsPerHost     int = 100
)

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func CustomerSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{
		Director:  director,
		Transport: DefaultTransport,
	}
}

var DefaultTransport http.RoundTripper = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: time.Duration(defaultTimeOut) * time.Second,
	}).Dial,
	TLSHandshakeTimeout: time.Duration(defaultTimeOut) * time.Second,
	MaxIdleConns:        maxIdleConns,
	IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
	DisableCompression:  true,
	MaxIdleConnsPerHost: maxIdleConnsPerHost,
	MaxConnsPerHost:     maxConnsPerHost,
}

// r.Any("/scan-api/*proxyPath", func(c *gin.Context) { Proxy("", "http://localhost:4002,http://localhost:4004", c) })
// r.Any("/scan-api/*proxyPath", func(c *gin.Context) { Proxy("/scan-api", "http://localhost:4002,http://localhost:4004", c) })
// r.Any("/a/scan-api/*proxyPath", func(c *gin.Context) { Proxy("/scan-api", "http://localhost:4002,http://localhost:4004", c) })
// r.Any("/b/scan-api/*proxyPath", func(c *gin.Context) { Proxy("/scan-api", "http://localhost:4002,http://localhost:4004", c) })
func Proxy(c *gin.Context, prefix string, remotes ...string) {

	rand.Seed(time.Now().UnixNano())

	upstreams := remotes
	upstreamSize := len(upstreams)
	currentRemote := upstreams[rand.Intn((upstreamSize-1)-0+1)+0]

	upstream, err := url.Parse(currentRemote)

	if err != nil {
		panic(err)
	}

	// httputil.ReverseProxy{}

	proxy := CustomerSingleHostReverseProxy(upstream)

	proxy.Director = func(req *http.Request) {

		RequestURI := req.URL.Path
		requestPath := strings.Replace((prefix + "" + c.Param("proxyPath")), "//", "/", -1)

		req.Header = c.Request.Header
		req.Host = upstream.Host
		req.URL.Scheme = upstream.Scheme
		req.URL.Host = upstream.Host
		req.URL.Path = requestPath

		// c.Request.WithContext(context.WithValue(c.Request.Context(), "ProxyReverse", utils.MapOf("Upstream", currentRemote, "RequestPath", requestPath)))

		logger.Infof(`Access-by-Proxy: URI=%s, Upstream=%s, ProxyPath=%s`, RequestURI, currentRemote, requestPath)
		// c.AbortWithStatus(201)

	}

	proxy.ServeHTTP(c.Writer, c.Request)

}

// gproxy.Any(r, "/a/scan-api", "/scan-api", "http://172.16.197.233:8080", "http://172.16.197.134:8080")
func Any(r *gin.Engine, from string, to string, remotes ...string) {
	r.Any(from+"/*proxyPath", func(c *gin.Context) {
		Proxy(c, to, remotes...)
	})
	logger.Infof(`Proxy-Reverse-Handler: Method=%s, From=%s, To=%s, remotes=%s`, "Any", from, to, strings.Join(remotes, ","))
}
