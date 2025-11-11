package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/app/router"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/handler"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/repository"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/service"
)

func main() {
	// Add logging
	log.Println("Starting Project Horizon Backend...")

	dbConnection, err := database.NewDatabaseConnection()
	if err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}
	defer dbConnection.Close()
	log.Println("Database connected successfully")

	analyticsRepository := repository.NewAnalyticsRepository(dbConnection)

	analyticsService := service.NewAnalyticsService(analyticsRepository)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Create Gin router
	r := gin.Default() // This includes logger & recovery middleware

	// Setup routes
	router.SetupRoutes(r, analyticsHandler)

	log.Println("Routes registered, starting server on :9095...")

	// Start server
	if err := r.Run(":9095"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
