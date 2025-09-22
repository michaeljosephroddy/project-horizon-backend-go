package analytics

import (
	"encoding/json"
	"fmt"

	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type AnalyticsService struct {
	journalRepository *database.JournalRepository
}

func NewAnalyticsService(journalRepository *database.JournalRepository) *AnalyticsService {
	return &AnalyticsService{
		journalRepository: journalRepository,
	}
}

func (service *AnalyticsService) Metrics(userID string, startDate string, endDate string) map[string]interface{} {

	movingAverages := service.journalRepository.MovingAverages(userID, startDate, endDate)

	trend := "constant"
	if movingAverages[len(movingAverages)-1].ThreeDay > movingAverages[len(movingAverages)-2].ThreeDay {
		trend = "increasing"
	}
	if movingAverages[len(movingAverages)-1].ThreeDay < movingAverages[len(movingAverages)-2].ThreeDay {
		trend = "decreasing"
	}
	fmt.Println(trend)

	standardDeviation := service.journalRepository.StandardDeviation(userID, startDate, endDate)

	stability := "stable"
	if standardDeviation >= 1.5 && standardDeviation < 3 {
		stability = "moderate"
	}
	if standardDeviation >= 3 {
		stability = "instable"
	}
	fmt.Println(stability)

	moodTagFrequencies := service.journalRepository.MoodTagFrequencies(userID, startDate, endDate)

	positiveDays := service.journalRepository.Days(userID, startDate, endDate, ">=", "6", "1", "50")

	negativeDays := service.journalRepository.Days(userID, startDate, endDate, "<=", "4", "2", "50")

	positiveStreaks := service.journalRepository.Streaks(userID, startDate, endDate, ">=", "6", "1", "50")

	negativeStreaks := service.journalRepository.Streaks(userID, startDate, endDate, "<=", "4", "2", "50")

	response := models.MetricsResponse{
		Trend:            trend,
		Stability:        stability,
		MoodTagFrequency: moodTagFrequencies,
		PositiveStreaks:  positiveStreaks,
		NegativeStreaks:  negativeStreaks,
		PositiveDays:     positiveDays,
		NegativeDays:     negativeDays,
	}

	fmt.Println("===================================")

	reponseJSON, _ := json.MarshalIndent(response, "", "    ")
	fmt.Println(string(reponseJSON))
	// highStreakCount := utils.StreakCount(highDays)
	// fmt.Println(highStreakCount)

	return map[string]interface{}{"empty": "empty"}
}
