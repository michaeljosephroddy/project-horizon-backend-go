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
	analyticsUsersMood       = `^/analytics/users/([0-9]+)/mood$`
	analyticsUsersSleep      = `^/analytics/users/([0-9]+)/sleep$`
	analyticsUsersMedication = `^/analytics/users/([0-9]+)/medication$`
)

func NewAnalyticsHandler(analyticsService *analyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (handler *AnalyticsHandler) ProcessRequest(writer http.ResponseWriter, request *http.Request) {
	switch {
	case common_utils.MatchURL(analyticsUsersMood, request.URL.Path):
		userID, startDate, endDate, err := handler.extractRequestParams(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		moodMetrics, err := handler.moodMetrics(userID, startDate, endDate)
		if err != nil {
			http.Error(writer, "failed to retrieve mood metrics", http.StatusInternalServerError)
			return
		}
		handler.writeJSONResponse(writer, moodMetrics)
	case common_utils.MatchURL(analyticsUsersSleep, request.URL.Path):
		userID, startDate, endDate, err := handler.extractRequestParams(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		sleepMetrics, err := handler.sleepMetrics(userID, startDate, endDate)
		if err != nil {
			http.Error(writer, "failed to retrieve sleep metrics", http.StatusInternalServerError)
			return
		}
		handler.writeJSONResponse(writer, sleepMetrics)
	default:
		http.Error(writer, "404 path not found", http.StatusNotFound)
	}
}
func (handler *AnalyticsHandler) extractRequestParams(request *http.Request) (string, time.Time, time.Time, error) {
	userID, err := common_utils.GetUserIDFromPath(request.URL.Path)
	if err != nil {
		return "", time.Time{}, time.Time{}, err
	}

	startDate := request.URL.Query().Get("startDate")
	endDate := request.URL.Query().Get("endDate")

	if startDate == "" || endDate == "" {
		return "", time.Time{}, time.Time{}, fmt.Errorf("startDate and endDate are required")
	}

	startDateParsed, endDateParsed := common_utils.ParseDates(startDate, endDate)

	return userID, startDateParsed, endDateParsed, nil
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

func (handler *AnalyticsHandler) moodMetrics(userID string, startDate time.Time, endDate time.Time) (*models.MoodMetric, error) {
	current, err := handler.analyticsService.analyzeMood(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze current mood period: %w", err)
	}

	previousStart, previousEnd := analytics_utils.PreviousDates(startDate, endDate)
	previous, err := handler.analyticsService.analyzeMood(userID, previousStart, previousEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze previous mood period: %w", err)
	}

	diffs := handler.analyticsService.moodDiffs(current, previous)
	current.Diffs = diffs

	return current, nil
}

func (handler *AnalyticsHandler) sleepMetrics(userID string, startDate time.Time, endDate time.Time) (*models.SleepMetric, error) {
	current, err := handler.analyticsService.analyzeSleep(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze current sleep period: %w", err)
	}

	previousStart, previousEnd := analytics_utils.PreviousDates(startDate, endDate)
	previous, err := handler.analyticsService.analyzeSleep(userID, previousStart, previousEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze previous sleep period: %w", err)
	}

	diffs := handler.analyticsService.sleepDiffs(current, previous)
	current.SleepDiffs = diffs

	return current, nil
}

// TODO implement this func
func (handler *AnalyticsHandler) medicationMetrics(userID string, startDate time.Time, endDate time.Time) (*models.MedicationMetric, error) {
	return &models.MedicationMetric{}, fmt.Errorf("todo")
}
