package analytics

import (
	"encoding/json"
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

func (service *AnalyticsService) Metrics(userID string, startDate string, endDate string) map[string]interface{} {

	movingAverages := service.journalRepository.MovingAverages(userID, startDate, endDate)
	movingAveragesJSON, _ := json.Marshal(movingAverages)
	fmt.Println("movingAverages ", string(movingAveragesJSON))
	// TODO implement increasing, decreasing or constant trend logic
	// if moving avg > yesterday avg ::: increasing 
	// if moving avg < yesterday avg ::: decreasing
	// else constant

	standardDeviation := service.journalRepository.StandardDeviation(userID, startDate, endDate)
	fmt.Println("standardDeviation ", standardDeviation)
	// TODO implement stability logic 
	// if std deviation < 1.5 ::: low volatility stable
	// if std deviation >= 1.5 and < 3 ::: moderate volatility
	// if std deviatiion >= 3 ::: high volatility instable 

	moodTagFrequencies := service.journalRepository.MoodTagFrequencies(userID, startDate, endDate)
	fmt.Println("moodTagFrequencies ", moodTagFrequencies)
	// TODO add to reponse object as a list of objects

	lowDays := service.journalRepository.LowDays(userID, startDate, endDate)
	fmt.Println("lowdays ", lowDays)
	// TODO impplement longest low streak
	// consecutive days where mood rating is 4 or below

	highDays := service.journalRepository.HighDays(userID, startDate, endDate)
	fmt.Println("highdays ", highDays)
	// TODO implement longest high steak
	// consecutive days where mood rating is 6 or above

	return map[string]interface{}{"empty": "empty"}
}
