package router

import (
	"github.com/gin-gonic/gin"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/handler"
)

func SetupRoutes(r *gin.Engine, analyticsHandler *handler.AnalyticsHandler) {
	// Analytics routes group
	analyticsGroup := r.Group("/analytics")
	{
		analyticsGroup.GET("/users/:userID/mood", analyticsHandler.GetMoodMetrics)
		analyticsGroup.GET("/users/:userID/sleep", analyticsHandler.GetSleepMetrics)
		analyticsGroup.GET("/users/:userID/medication", analyticsHandler.GetMedicationMetrics)
	}
}
