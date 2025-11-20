package router

import (
	"github.com/gin-gonic/gin"
	analyticshandler "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/handler"
	authhandler "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/handler"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/middleware"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/service"
	usershandler "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/handler"
)

func SetupRoutes(
	r *gin.Engine,
	analyticsHandler *analyticshandler.AnalyticsHandler,
	authHandler *authhandler.AuthHandler,
	authService *service.AuthService,
	usersHandler *usershandler.UserHandler,
) {
	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// Auth endpoints
		protected.GET("/auth/profile", authHandler.GetProfile)

		// Analytics endpoints (protected)
		analytics := protected.Group("/analytics")
		{
			analytics.GET("/users/:userID/mood", analyticsHandler.GetMoodMetrics)
			analytics.GET("/users/:userID/sleep", analyticsHandler.GetSleepMetrics)
			analytics.GET("/users/:userID/medication", analyticsHandler.GetMedicationMetrics)
		}

		users := protected.Group("/users")
		{
			users.GET("/:userID/user-medication", usersHandler.GetUserMedications)
		}

		// Other protected endpoints
		// protected.GET("/users", getUsersHandler)
	}
}
