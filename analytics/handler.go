package analytics

import (
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
	analytics_utils "github.com/michaeljosephroddy/project-horizon-backend-go/analytics/utils"
	common_utils "github.com/michaeljosephroddy/project-horizon-backend-go/common/utils"
)

type AnalyticsHandler struct {
	analyticsService *analyticsService
}

func NewAnalyticsHandler(analyticsService *analyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetMoodMetrics handles GET /analytics/users/:userID/mood
func (handler *AnalyticsHandler) GetMoodMetrics(c *gin.Context) {
	userID, startDate, endDate, err := handler.extractRequestParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	moodMetrics, err := handler.moodMetrics(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve mood metrics"})
		return
	}

	c.JSON(http.StatusOK, moodMetrics)
}

// GetSleepMetrics handles GET /analytics/users/:userID/sleep
func (handler *AnalyticsHandler) GetSleepMetrics(c *gin.Context) {
	userID, startDate, endDate, err := handler.extractRequestParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sleepMetrics, err := handler.sleepMetrics(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve sleep metrics"})
		return
	}

	c.JSON(http.StatusOK, sleepMetrics)
}

// GetMedicationMetrics handles GET /analytics/users/:userID/medication
func (handler *AnalyticsHandler) GetMedicationMetrics(c *gin.Context) {
	userID, startDate, endDate, err := handler.extractRequestParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medicationMetrics, err := handler.medicationMetrics(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve medication metrics"})
		return
	}

	c.JSON(http.StatusOK, medicationMetrics)
}

func (handler *AnalyticsHandler) extractRequestParams(c *gin.Context) (string, time.Time, time.Time, error) {
	userID := c.Param("userID")
	if userID == "" {
		return "", time.Time{}, time.Time{}, fmt.Errorf("userID is required")
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	if startDate == "" || endDate == "" {
		return "", time.Time{}, time.Time{}, fmt.Errorf("startDate and endDate are required")
	}

	startDateParsed, endDateParsed := common_utils.ParseDates(startDate, endDate)

	return userID, startDateParsed, endDateParsed, nil
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

func (handler *AnalyticsHandler) medicationMetrics(userID string, startDate time.Time, endDate time.Time) (*models.MedicationMetric, error) {
	current, err := handler.analyticsService.analyzeMedication(userID, startDate, endDate)
	if err != nil {
		log.Println("somehting wrong here %w", err)	
		return nil, fmt.Errorf("failed to analyze current medication period: %w", err)
	}

	previousStart, previousEnd := analytics_utils.PreviousDates(startDate, endDate)
	previous, err := handler.analyticsService.analyzeMedication(userID, previousStart, previousEnd)
	if err != nil {
		log.Println("somehting wrong here as well %w", err)	
		return nil, fmt.Errorf("failed to analyze previous medication period: %w", err)
	}

	diffs := handler.analyticsService.medicationDiffs(current, previous)
	current.MedicationDiffs = diffs

	return current, nil
}
