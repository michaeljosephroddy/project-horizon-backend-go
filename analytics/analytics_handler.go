package analytics

import (
	"encoding/json"
	"fmt"
	"github.com/michaeljosephroddy/project-horizon-backend-go/utils"
	"net/http"
)

type AnalyticsHandler struct {
	analyticsService *analyticsService
}

// TODO need to come up with a better regexp
var metricsRegexp string = `^/analytics/users/([0-9]+)/metrics$`

func NewAnalyticsHandler(analyticsService *analyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (handler *AnalyticsHandler) ProcessRequest(writer http.ResponseWriter, request *http.Request) {
	switch {
	case utils.MatchURL(metricsRegexp, request.URL.Path):

		userID := utils.GetUserIDFromPath(request.URL.Path)
		startDate := request.URL.Query().Get("startDate")
		endDate := request.URL.Query().Get("endDate")

		result := handler.analyticsService.metrics(userID, startDate, endDate)
		fmt.Printf("DEBUG struct: %+v\n", result)
		body, _ := json.Marshal(result)

		fmt.Println("DEBUG json:", string(body))

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
