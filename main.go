package main

import (
	"net/http"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics"
	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/router"
)

func main() {

	dbConnection := database.NewDatabaseConnection()
	defer dbConnection.Close()

	moodLogRepository := database.NewMoodLogRepository(dbConnection)
	sleepLogRepository := database.NewSleepLogRepository(dbConnection)
	analyticsService := analytics.NewAnalyticsService(moodLogRepository, sleepLogRepository)
	analyticsHandler := analytics.NewAnalyticsHandler(analyticsService)
	r := router.NewRouter(analyticsHandler)

	http.HandleFunc("/", r.RouteRequests)
	http.ListenAndServe(":9095", nil)
}
