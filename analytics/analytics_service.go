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

	// movingAveragesJSON, _ := json.MarshalIndent(movingAverages, "", "    ")
	// fmt.Println(string(movingAveragesJSON))

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
	moodTagFrequenciesJSON, _ := json.MarshalIndent(moodTagFrequencies, "", "    ")
	fmt.Println(string(moodTagFrequenciesJSON))

	// TODO fix this problem with >= 6 and <= 4 for sql
	// i should be able to use the same func Days for both conditions
	// make it work
	highDays := service.journalRepository.Days(userID, startDate, endDate, ">=", "6")
	highDaysJSON, _ := json.MarshalIndent(highDays, "", "    ")
	fmt.Println(string(highDaysJSON))

	lowDays := service.journalRepository.Days(userID, startDate, endDate, "<=", "4")
	lowDaysJSON, _ := json.MarshalIndent(lowDays, "", "    ")
	fmt.Println(string(lowDaysJSON))

	// TODO same goes for Streaks
	// fix this problem with >= 6 and <= 4 for sql
	// I should be able to use the same func Streaks for both conditions
	// make it work
	highStreaks := service.journalRepository.Streaks(userID, startDate, endDate, ">=", "6")
	highStreaksJSON, _ := json.MarshalIndent(highStreaks, "", "    ")
	fmt.Println(string(highStreaksJSON))

	lowStreaks := service.journalRepository.Streaks(userID, startDate, endDate, "<=", "4")
	lowStreaksJSON, _ := json.MarshalIndent(lowStreaks, "", "    ")
	fmt.Println(string(lowStreaksJSON))

	response := models.MetricsResponse{
		Trend: trend,
		Stability: stability,
		MoodTagFrequency: moodTagFrequencies,
		HighStreaks: highStreaks,
		LowStreaks: lowStreaks,
	}

	fmt.Println("===================================")

	reponseJSON, _ := json.MarshalIndent(response, "", "    ")
	fmt.Println(string(reponseJSON))
	// highStreakCount := utils.StreakCount(highDays)
	// fmt.Println(highStreakCount)

	return map[string]interface{}{"empty": "empty"}
}
