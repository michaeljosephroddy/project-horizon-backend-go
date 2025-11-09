package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics"
	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/router"
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

	moodLogRepository := database.NewMoodLogRepository(dbConnection)
	sleepLogRepository := database.NewSleepLogRepository(dbConnection)
	medicationLogRepository := database.NewMedicationLogRepository(dbConnection)

	analyticsService := analytics.NewAnalyticsService(moodLogRepository, sleepLogRepository, medicationLogRepository)
	analyticsHandler := analytics.NewAnalyticsHandler(analyticsService)

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
