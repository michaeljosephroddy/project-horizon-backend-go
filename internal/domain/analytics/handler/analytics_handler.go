package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/common/utils"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/service"
	authutils "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/utils"
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
func (ah *AnalyticsHandler) GetMoodMetrics(c *gin.Context) {
	userID, startDate, endDate, err := ah.extractRequestParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authorization check - ensure user can only access their own data
	if !authutils.AuthorizeUserAccess(c, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	moodMetrics, err := ah.analyticsService.AnalyzeMood(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve mood metrics"})
		return
	}

	c.JSON(http.StatusOK, moodMetrics)
}

// GetSleepMetrics handles GET /analytics/users/:userID/sleep
func (ah *AnalyticsHandler) GetSleepMetrics(c *gin.Context) {
	userID, startDate, endDate, err := ah.extractRequestParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authorization check
	if !authutils.AuthorizeUserAccess(c, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	sleepMetrics, err := ah.analyticsService.AnalyzeSleep(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve sleep metrics"})
		return
	}

	c.JSON(http.StatusOK, sleepMetrics)
}

// GetMedicationMetrics handles GET /analytics/users/:userID/medication
func (ah *AnalyticsHandler) GetMedicationMetrics(c *gin.Context) {
	userID, startDate, endDate, err := ah.extractRequestParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authorization check
	if !authutils.AuthorizeUserAccess(c, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	medicationMetrics, err := ah.analyticsService.AnalyzeMedication(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve medication metrics"})
		return
	}

	c.JSON(http.StatusOK, medicationMetrics)
}

func (ah *AnalyticsHandler) extractRequestParams(c *gin.Context) (int, time.Time, time.Time, error) {
	userID, _ := strconv.Atoi(c.Param("userID"))
	if userID == 0 {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("userID is required")
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	if startDate == "" || endDate == "" {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("startDate and endDate are required")
	}

	startDateParsed := utils.ParseDate(startDate)
	endDateParsed := utils.ParseDate(endDate)

	return userID, startDateParsed, endDateParsed, nil
}
