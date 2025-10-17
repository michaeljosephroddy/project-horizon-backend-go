package analytics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
	analytics_utils "github.com/michaeljosephroddy/project-horizon-backend-go/analytics/utils"
	common_utils "github.com/michaeljosephroddy/project-horizon-backend-go/common/utils"
)

type AnalyticsHandler struct {
	analyticsService *analyticsService
}

const (
	analyticsUsersMood  = `^/analytics/users/([0-9]+)/mood$`
	analyticsUsersSleep = `^/analytics/users/([0-9]+)/sleep$`
	// analyticsUsersMedication = `^/analytics/users/([0-9]+)/medication$`
)

func NewAnalyticsHandler(analyticsService *analyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (handler *AnalyticsHandler) ProcessRequest(writer http.ResponseWriter, request *http.Request) {
	switch {
	case common_utils.MatchURL(analyticsUsersMood, request.URL.Path):
		handler.handleMoodRequest(writer, request)
	case common_utils.MatchURL(analyticsUsersSleep, request.URL.Path):
		handler.handleSleepRequest(writer, request)
	default:
		http.Error(writer, "404 path not found", http.StatusNotFound)
	}
}

func (handler *AnalyticsHandler) handleMoodRequest(writer http.ResponseWriter, request *http.Request) {
	userID, startDate, endDate, err := handler.extractRequestParams(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	moodMetrics := handler.moodMetrics(userID, startDate, endDate)
	handler.writeJSONResponse(writer, moodMetrics)
}

func (handler *AnalyticsHandler) handleSleepRequest(writer http.ResponseWriter, request *http.Request) {
	userID, startDate, endDate, err := handler.extractRequestParams(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	sleepMetrics := handler.sleepMetrics(userID, startDate, endDate)
	handler.writeJSONResponse(writer, sleepMetrics)
}

func (handler *AnalyticsHandler) extractRequestParams(request *http.Request) (string, time.Time, time.Time, error) {
	userID := common_utils.GetUserIDFromPath(request.URL.Path)
	if userID == "" {
		return "", time.Time{}, time.Time{}, fmt.Errorf("invalid user ID")
	}

	startDateStr := request.URL.Query().Get("startDate")
	endDateStr := request.URL.Query().Get("endDate")
	
	if startDateStr == "" || endDateStr == "" {
		return "", time.Time{}, time.Time{}, fmt.Errorf("startDate and endDate are required")
	}

	startDate, endDate := common_utils.ParseDates(startDateStr, endDateStr)
	return userID, startDate, endDate, nil
}

func (handler *AnalyticsHandler) writeJSONResponse(writer http.ResponseWriter, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(body)
}

func (handler *AnalyticsHandler) moodMetrics(userID string, startDate time.Time, endDate time.Time) *models.MoodMetric {
	current := handler.analyticsService.analyzeMood(userID, startDate, endDate)
	previousStart, previousEnd := analytics_utils.PreviousDates(startDate, endDate)
	previous := handler.analyticsService.analyzeMood(userID, previousStart, previousEnd)
	diffs := handler.analyticsService.moodDiffs(current, previous)
	current.Diffs = diffs
	return current
}

func (handler *AnalyticsHandler) sleepMetrics(userID string, startDate time.Time, endDate time.Time) *models.SleepMetric {
	current := handler.analyticsService.analyzeSleep(userID, startDate, endDate)
	previousStart, previousEnd := analytics_utils.PreviousDates(startDate, endDate)
	previous := handler.analyticsService.analyzeSleep(userID, previousStart, previousEnd)
	diffs := handler.analyticsService.sleepDiffs(current, previous)
	current.SleepDiffs = diffs
	return current
}
