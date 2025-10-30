package router

import (
	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics"
	"net/http"
	"strings"
)

type Router struct {
	analyticsHandler *analytics.AnalyticsHandler
}

func NewRouter(handler *analytics.AnalyticsHandler) *Router {
	return &Router{analyticsHandler: handler}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch {
	case strings.HasPrefix(req.URL.Path, "/analytics"):
		r.analyticsHandler.ProcessRequest(w, req)
	default:
		http.NotFound(w, req)
	}
}
