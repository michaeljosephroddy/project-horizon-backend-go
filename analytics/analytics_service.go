package analytics

import (
	"slices"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/utils"

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

func (service *AnalyticsService) Metrics(userID string, startDate string, endDate string) models.Metrics {

	currentPeriodMetrics := service.periodAnalytics(userID, startDate, endDate)

	// TODO calcuate previous period start and end dates
	// the previous period length should be == to the current period length
	// the previous period start date should be current period start - period length
	// Assume today is the end of the current period

	// TODO take a look at the moving averges trend implementation
	// it should cover the period not just 3 days
	// maybe daily is 3 days, weekly 7 days and monthly 28 days moving avg

	layout := "2006-01-02"
	startDateParsed, _ := time.Parse(layout, startDate)
	endDateParsed, _ := time.Parse(layout, endDate)
	diff := endDateParsed.Sub(startDateParsed)
	numDays := int(diff.Hours() / 24)
	previousStart := startDateParsed.AddDate(0, 0, -numDays)
	previousEnd := startDateParsed.AddDate(0, 0, -1)

	previousPeriodMetrics := service.periodAnalytics(userID, previousStart.Format(layout), previousEnd.Format(layout))

	periodDiffs := service.diffs(currentPeriodMetrics, previousPeriodMetrics)

	currentPeriodMetrics.PeriodDiffs = periodDiffs

	return currentPeriodMetrics
}
func (service *AnalyticsService) periodAnalytics(userID string, startDate string, endDate string) models.Metrics {

	movingAverages := service.journalRepository.MovingAverages(userID, startDate, endDate)

	var trend string

	if len(movingAverages) >= 2 {
		last := movingAverages[len(movingAverages)-1]
		prev := movingAverages[len(movingAverages)-2]

		switch {
		case last.ThreeDay > prev.ThreeDay:
			trend = "increasing"
		case last.ThreeDay < prev.ThreeDay:
			trend = "decreasing"
		default:
			trend = "flat"
		}
	} else {
		trend = "not enough data"
	}

	standardDeviation := service.journalRepository.StandardDeviation(userID, startDate, endDate)

	var stability string

	switch {
	case standardDeviation == 0:
		stability = "not enough data" // e.g., only 1 data point
	case standardDeviation < 1.5:
		stability = "stable"
	case standardDeviation < 3:
		stability = "moderate"
	default:
		stability = "unstable"
	}

	avgMoodRatingPeriod := service.journalRepository.PeriodAvgMoodRating(userID, startDate, endDate)

	mtfPeriod := service.journalRepository.MoodTagFrequencies(userID, startDate, endDate)

	slices.SortFunc(mtfPeriod, func(a, b models.MoodTagFrequency) int {
		if a.Percentage > b.Percentage {
			return -1
		} else if a.Percentage < b.Percentage {
			return 1
		} else {
			return 0
		}
	})

	positiveDays := service.journalRepository.Days(userID, startDate, endDate, ">=", "6", "1", "50")

	mtfPositiveDays := utils.MoodTagFrequencies(positiveDays)

	// TODO find top mood tags for positive days
	// i.e Happy is the most frequently recorded tag on good days
	// i.e HAPPY and CONTENT are the most frequently recorded tags on good days
	// i.e HAPPY, CONTENT and CALM are the most frequently recorded tags on good days

	negativeDays := service.journalRepository.Days(userID, startDate, endDate, "<=", "4", "2", "50")

	mtfNegativeDays := utils.MoodTagFrequencies(negativeDays)

	positiveStreaks := service.journalRepository.Streaks(userID, startDate, endDate, ">=", "6", "1", "50")

	negativeStreaks := service.journalRepository.Streaks(userID, startDate, endDate, "<=", "4", "2", "50")

	layout := "2006-01-02" // Correct Go layout
	startDateParsed, _ := time.Parse(layout, startDate)
	endDateParsed, _ := time.Parse(layout, endDate)

	diff := endDateParsed.Sub(startDateParsed)
	numDays := int(diff.Hours() / 24)

	var granularity string
	switch {
	case numDays <= 1:
		granularity = "daily"
	case numDays <= 7:
		granularity = "weekly"
	case numDays <= 28:
		granularity = "monthly"
	default:
		granularity = "custom"
	}

	metrics := models.Metrics{
		UserID:               userID,
		Granularity:          granularity,
		PeriodStart:          startDate,
		PeriodEnd:            endDate,
		Trend:                trend,
		Stability:            stability,
		AvgMoodRatingPeriod:  avgMoodRatingPeriod,
		TopMoodsPeriod:       mtfPeriod,
		TopMoodsPositiveDays: mtfPositiveDays,
		TopMoodsNegativeDays: mtfNegativeDays,
		PositiveStreaks:      positiveStreaks,
		NegativeStreaks:      negativeStreaks,
		PositiveDays:         positiveDays,
		NegativeDays:         negativeDays,
		PeriodDiffs:          models.PeriodDiff{},
	}

	return metrics
}

func (servcie *AnalyticsService) diffs(currentPeriod models.Metrics, previousPeriod models.Metrics) models.PeriodDiff {

	return models.PeriodDiff{}

}
