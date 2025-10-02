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
// var analyticsUsersMedication string = `^/analytics/users/([0-9]+)/medication$`

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

		moodMetrics := handler.moodMetrics(userID, startDate, endDate)
		body, _ := json.Marshal(moodMetrics)
		fmt.Println("DEBUG ", string(body))

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(body)

	case utils.MatchURL(analyticsUsersSleep, request.URL.Path):

		userID := utils.GetUserIDFromPath(request.URL.Path)
		startDate := request.URL.Query().Get("startDate")
		endDate := request.URL.Query().Get("endDate")

		sleepMetrics := handler.sleepMetrics(userID, startDate, endDate)
		body, _ := json.Marshal(sleepMetrics)
		fmt.Println("DEBUG ", string(body))

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(body)

	/* case utils.MatchURL(analyticsUsersMedication, request.URL.Path):

	userID := utils.GetUserIDFromPath(request.URL.Path)
	startDate := request.URL.Query().Get("startDate")
	endDate := request.URL.Query().Get("endDate")

	var medication models.Medication
	body, _ := json.Marshal(medication)
	fmt.Println("DEBUG ", string(body))

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(body) */
	default:
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("404 path not found"))
	}
}

func (handler *AnalyticsHandler) moodMetrics(userID string, startDate string, endDate string) *models.MoodMetric {

	current := handler.analyticsService.analyzeMood(userID, startDate, endDate)
	fmt.Println(startDate, " ", endDate, current.MovingAvg)

	previousStart, previousEnd := utils.CalculatePreviousDates(startDate, endDate)

	previous := handler.analyticsService.analyzeMood(userID, previousStart, previousEnd)
	fmt.Println(previousStart, " ", previousEnd, previous.MovingAvg)

	diffs := handler.analyticsService.moodDiffs(current, previous)

	current.MoodDiffs = diffs

	return current
}

func (handler *AnalyticsHandler) sleepMetrics(userID string, startDate string, endDate string) *models.SleepMetric {

	current := handler.analyticsService.analyzeSleep(userID, startDate, endDate)

	return current 
}
