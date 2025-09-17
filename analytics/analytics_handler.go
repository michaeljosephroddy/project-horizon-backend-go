package analytics

import (
	"encoding/json"
	"github.com/michaeljosephroddy/project-horizon-backend-go/utils"
	"net/http"
)

type AnalyticsHandler struct {
	AnalyticsService *AnalyticsService
}

var metricsRegexp string = `^/analytics/users/([0-9]+)/metrics$`

func NewAnalyticsHandler(analyticsService *AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		AnalyticsService: analyticsService,
	}
}

func (handler *AnalyticsHandler) ProcessRequest(writer http.ResponseWriter, request *http.Request) {

	switch {
	case utils.MatchURL(metricsRegexp, request.URL.Path):

		userId := utils.GetUserIdFromPath(request.URL.Path)
		startDate := request.URL.Query().Get("startDate")
		endDate := request.URL.Query().Get("endDate")

		result := handler.AnalyticsService.Metrics(userId, startDate, endDate)
		body, _ := json.Marshal(result)

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(body)
	// case patternMatch("/api/analytics/users/(\\d+)/period-comparison", request.URL.Path):
	// 	processPeriodComparison(request, writer)
	// case patternMatch("/api/analytics/users/(\\d+)/trend-analysis", request.URL.Path):
	// 	processTrendAnalysis(request, writer)
	default:
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("404 path not found"))
	}
}
