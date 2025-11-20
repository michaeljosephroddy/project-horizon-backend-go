package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/app/router"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/database"

	analyticshandler "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/handler"
	analyticsrepository "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/repository"
	analyticsservice "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/service"
	authhandler "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/handler"
	authservice "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/service"
	usershandler "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/handler"
	usersrepository "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/repository"
	usersservice "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/service"
)

func main() {
	// Add logging
	log.Println("Starting Project Horizon Backend...")

	// Initialize database connection
	dbConnection, err := database.NewDatabaseConnection()
	if err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}
	defer dbConnection.Close()
	log.Println("Database connected successfully")

	// Initialize Analytics domain
	analyticsRepository := analyticsrepository.NewAnalyticsRepository(dbConnection)
	analyticsService := analyticsservice.NewAnalyticsService(analyticsRepository)
	analyticsHandler := analyticshandler.NewAnalyticsHandler(analyticsService)

	usersRepository := usersrepository.NewUsersRepository(dbConnection)
	usersService := usersservice.NewUsersService(usersRepository)
	usersHandler := usershandler.NewUsersHandler(usersService)

	// Initialize Auth domainb
	authService := authservice.NewAuthService(usersRepository)
	authHandler := authhandler.NewAuthHandler(authService)

	// Create Gin router
	r := gin.Default() // This includes logger & recovery middleware

	// Setup routes (pass all required handlers and services)
	router.SetupRoutes(r, analyticsHandler, authHandler, authService, usersHandler)

	log.Println("Routes registered, starting server on :9095...")

	// Start server
	if err := r.Run(":9095"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
