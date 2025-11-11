package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	common_utils "github.com/michaeljosephroddy/project-horizon-backend-go/internal/common/utils"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/service"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
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

	moodMetrics, err := handler.analyticsService.AnalyzeMood(userID, startDate, endDate)
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

	sleepMetrics, err := handler.analyticsService.AnalyzeSleep(userID, startDate, endDate)
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

	medicationMetrics, err := handler.analyticsService.AnalyzeMedication(userID, startDate, endDate)
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
