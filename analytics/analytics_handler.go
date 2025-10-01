package analytics

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
	"github.com/michaeljosephroddy/project-horizon-backend-go/utils"
)

type AnalyticsHandler struct {
	analyticsService *analyticsService
}

// TODO need to come up with a better regexp
var analyticsUsersMood string = `^/analytics/users/([0-9]+)/mood$`
var analyticsUsersSleep string = `^/analytics/users/([0-9]+)/sleep$`
var analyticsUsersMedication string = `^/analytics/users/([0-9]+)/medication$`

func NewAnalyticsHandler(analyticsService *analyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (handler *AnalyticsHandler) ProcessRequest(writer http.ResponseWriter, request *http.Request) {
	switch {
	case utils.MatchURL(analyticsUsersMood, request.URL.Path):

		userID := utils.GetUserIDFromPath(request.URL.Path)
		startDate := request.URL.Query().Get("startDate")
		endDate := request.URL.Query().Get("endDate")

		result := handler.analyticsService.metrics(userID, startDate, endDate)
		body, _ := json.Marshal(result)
		fmt.Println("DEBUG ", string(body))

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(body)

	case utils.MatchURL(analyticsUsersSleep, request.URL.Path):

		/* userID := utils.GetUserIDFromPath(request.URL.Path)
		startDate := request.URL.Query().Get("startDate")
		endDate := request.URL.Query().Get("endDate") */

		var sleep models.Sleep
		body, _ := json.Marshal(sleep)
		fmt.Println("DEBUG ", string(body))

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(body)

	case utils.MatchURL(analyticsUsersMedication, request.URL.Path):

		/* userID := utils.GetUserIDFromPath(request.URL.Path)
		startDate := request.URL.Query().Get("startDate")
		endDate := request.URL.Query().Get("endDate") */

		var medication models.Medication
		body, _ := json.Marshal(medication)
		fmt.Println("DEBUG ", string(body))

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(body)

	// 	processPeriodComparison(request, writer)
	// case patternMatch("/api/analytics/users/(\\d+)/trend-analysis", request.URL.Path):
	// 	processTrendAnalysis(request, writer)
	default:
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("404 path not found"))
	}
}
