package analytics

import (
	"slices"

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

func (service *AnalyticsService) Metrics(userID string, startDate string, endDate string) models.MetricsResponse {

	movingAverages := service.journalRepository.MovingAverages(userID, startDate, endDate)

	trend := "constant"
	if movingAverages[len(movingAverages)-1].ThreeDay > movingAverages[len(movingAverages)-2].ThreeDay {
		trend = "increasing"
	}
	if movingAverages[len(movingAverages)-1].ThreeDay < movingAverages[len(movingAverages)-2].ThreeDay {
		trend = "decreasing"
	}

	standardDeviation := service.journalRepository.StandardDeviation(userID, startDate, endDate)

	stability := "stable"
	if standardDeviation >= 1.5 && standardDeviation < 3 {
		stability = "moderate"
	}
	if standardDeviation >= 3 {
		stability = "instable"
	}

	moodTagFrequencies := service.journalRepository.MoodTagFrequencies(userID, startDate, endDate)

	slices.SortFunc(moodTagFrequencies, func(a, b models.MoodTagFrequency) int {
		if a.Percentage > b.Percentage {
			return -1
		} else if a.Percentage < b.Percentage {
			return 1
		} else {
			return 0
		}
	})

	top3Moods := moodTagFrequencies[:3]

	positiveDays := service.journalRepository.Days(userID, startDate, endDate, ">=", "6", "1", "50")

	// mtfPositiveDays := utils.MoodTagFrequencies(positiveDays)

	// TODO find top mood tags for positive days
	// i.e Happy is the most frequently recorded tag on good days
	// i.e HAPPY and CONTENT are the most frequently recorded tags on good days
	// i.e HAPPY, CONTENT and CALM are the most frequently recorded tags on good days

	negativeDays := service.journalRepository.Days(userID, startDate, endDate, "<=", "4", "2", "50")

	positiveStreaks := service.journalRepository.Streaks(userID, startDate, endDate, ">=", "6", "1", "50")

	negativeStreaks := service.journalRepository.Streaks(userID, startDate, endDate, "<=", "4", "2", "50")

	response := models.MetricsResponse{
		UserID:             userID,
		PeriodStart:        startDate,
		PeriodEnd:          endDate,
		Trend:              trend,
		Stability:          stability,
		MoodTagFrequencies: moodTagFrequencies,
		Top3Moods:          top3Moods,
		PositiveStreaks:    positiveStreaks,
		NegativeStreaks:    negativeStreaks,
		PositiveDays:       positiveDays,
		NegativeDays:       negativeDays,
	}

	// reponseJSON, _ := json.MarshalIndent(response, "", "    ")
	// fmt.Println(string(reponseJSON))

	return response
}
