package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ProxyHandler creates a reverse proxy handler for the given target base URL.
// Preserves the original request path, query, and all headers (including X-User-ID, X-User-Role).
func ProxyHandler(targetBase string) gin.HandlerFunc {
	targetURL, _ := url.Parse(targetBase)

	return func(c *gin.Context) {
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = targetURL.Scheme
				req.URL.Host = targetURL.Host
				req.Host = targetURL.Host
				// Path, query and headers are already set by Gin / JWT middleware
			},
		}
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
