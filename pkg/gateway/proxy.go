package gateway

import (
	"github.com/gin-gonic/gin"
	"neotype-backend/pkg/config"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Proxy(c *gin.Context, service, endpoint string) {
	path, err := target(service, endpoint)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	targetUrl, err := url.Parse(path)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	proxy(targetUrl).ServeHTTP(c.Writer, c.Request)
}

// target returns a full URL to a service, determining the host from config (found by 'service' parameter)
// Returns something like: service/endpoint, so localhost:5010/endpoint
func target(service, endpoint string) (string, error) {
	host, err := config.GetBaseURL(service)
	if err != nil {
		return "", err
	}

	host += endpoint

	return host, nil
}

func proxy(address *url.URL) *httputil.ReverseProxy {
	p := httputil.NewSingleHostReverseProxy(address)

	p.Director = func(request *http.Request) {
		request.Host = address.Host
		request.URL.Scheme = address.Scheme
		request.URL.Host = address.Host
		request.URL.Path = address.Path
	}

	return p
}
