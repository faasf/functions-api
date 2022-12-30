package v1

import (
	"github.com/faasf/functions-api/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Proxy(cfg *config.Config, c *gin.Context) {
	remote, err := url.Parse(cfg.HTTP.NodeJsRuntimeUrl)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = "/call"
		req.Method = "POST"
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
