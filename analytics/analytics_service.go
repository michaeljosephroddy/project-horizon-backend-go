package analytics

import (
	"fmt"
	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
)

type AnalyticsService struct {
	journalRepository *database.JournalRepository
}

func NewAnalyticsService(journalRepository *database.JournalRepository) *AnalyticsService {
	return &AnalyticsService{
		journalRepository: journalRepository,
	}
}

func (service *AnalyticsService) Metrics(userId string, startDate string, endDate string) map[string]interface{} {

	movingAverages := service.journalRepository.MovingAverages(userId, startDate, endDate)
	fmt.Println("movingAverages ", movingAverages)

	standardDeviation := service.journalRepository.StandardDeviation(userId, startDate, endDate)
	fmt.Println("standardDeviation ", standardDeviation)

	moodTagFrequencies := service.journalRepository.MoodTagFrequencies(userId, startDate, endDate)
	fmt.Println("moodTagFrequencies ", moodTagFrequencies)

	lowDays := service.journalRepository.LowDays(userId, startDate, endDate)
	fmt.Println("lowdays ", lowDays)

	highDays := service.journalRepository.HighDays(userId, startDate, endDate)
	fmt.Println("highdays ", highDays)

	return map[string]interface{}{"empty": "empty"}
}
