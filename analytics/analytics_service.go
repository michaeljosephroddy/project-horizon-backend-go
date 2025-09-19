package analytics

import (
	"encoding/json"
	"fmt"
	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/utils"
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
	// TODO add to reponse object as a list of objects

	lowDays := service.journalRepository.LowDays(userID, startDate, endDate)
	lowDaysJSON, _ := json.MarshalIndent(lowDays, "", "    ")
	fmt.Println(string(lowDaysJSON))
	// TODO impplement longest low streak
	// consecutive days where mood rating is 4 or below
	// lowStreakCount := utils.StreakCount(lowDays)
	// fmt.Println(lowStreakCount)

	highDays := service.journalRepository.HighDays(userID, startDate, endDate)
	highDaysJSON, _ := json.MarshalIndent(highDays, "", "    ")
	fmt.Println(string(highDaysJSON))

	// highStreakCount := utils.StreakCount(highDays)
	// fmt.Println(highStreakCount)

	return map[string]interface{}{"empty": "empty"}
}
