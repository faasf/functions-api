package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Proxy(c *gin.Context) {
	remote, err := url.Parse("http://localhost:8082")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = "/test"
		req.Method = "POST"
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
