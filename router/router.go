package router

import (
	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics"
	"net/http"
	"strings"
)

type Router struct {
	AnalyticsHandler *analytics.AnalyticsHandler
}

func NewRouter(handler *analytics.AnalyticsHandler) *Router {
	return &Router{
		AnalyticsHandler: handler,
	}
}

func (r *Router) RouteRequests(writer http.ResponseWriter, request *http.Request) {
	switch {
	case strings.HasPrefix(request.URL.Path, "/analytics"):
		r.AnalyticsHandler.ProcessRequest(writer, request)
	default:
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("resouce not found"))
	}
}
