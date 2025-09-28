package analytics

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/utils"

	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type analyticsService struct {
	journalRepository *database.JournalRepository
}

func NewAnalyticsService(journalRepository *database.JournalRepository) *analyticsService {
	return &analyticsService{
		journalRepository: journalRepository,
	}
}

func (service *analyticsService) metrics(userID string, startDate string, endDate string) models.Metrics {

	current := service.analyze(userID, startDate, endDate)

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
	previousStart := startDateParsed.AddDate(0, 0, -numDays).Format(layout)
	previousEnd := startDateParsed.AddDate(0, 0, -1).Format(layout)

	previous := service.analyze(userID, previousStart, previousEnd)

	diffs := service.diffs(current, previous)

	current.Diffs = diffs

	return current
}

func (service *analyticsService) analyze(userID string, startDate string, endDate string) models.Metrics {

	movingAverages := service.journalRepository.MovingAverages(userID, startDate, endDate)

	var trend string

	if len(movingAverages) >= 2 {
		last := movingAverages[len(movingAverages)-1]
		prev := movingAverages[len(movingAverages)-2]

		switch {
		case last.MovingAvg > prev.MovingAvg:
			trend = "increasing"
		case last.MovingAvg < prev.MovingAvg:
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

	avgMoodRating := service.journalRepository.AvgMoodRating(userID, startDate, endDate)

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
		StartDate:            startDate,
		EndDate:              endDate,
		Trend:                trend,
		StdDeviation:         standardDeviation,
		Stability:            stability,
		AvgMoodRating:        avgMoodRating,
		TopMoods:             mtfPeriod,
		TopMoodsPositiveDays: mtfPositiveDays,
		TopMoodsNegativeDays: mtfNegativeDays,
		PositiveStreaks:      positiveStreaks,
		NegativeStreaks:      negativeStreaks,
		PositiveDays:         positiveDays,
		NegativeDays:         negativeDays,
		Diffs:                models.Diff{},
	}

	return metrics
}

func (servcie *analyticsService) diffs(current models.Metrics, previous models.Metrics) models.Diff {
	// [(New Value - Original Value) / Original Value] × 100

	var avgMoodChange float64
	if previous.AvgMoodRating != 0.0 {
		avgMoodChange = ((current.AvgMoodRating - previous.AvgMoodRating) / previous.AvgMoodRating) * 100
	}

	trendChange := fmt.Sprintf("%s -> %s", previous.Trend, current.Trend)

	stabilityChange := fmt.Sprintf("%s -> %s", previous.Stability, current.Stability)

	volatilityDelta := current.StdDeviation - previous.StdDeviation

	var volatilityDeltaPercentage float64
	if previous.StdDeviation != 0.0 {
		volatilityDeltaPercentage = ((current.StdDeviation - previous.StdDeviation) / previous.StdDeviation) * 100
	}

	// TODO wirte utility method
	var topMoodShift string
	if len(previous.TopMoods) >= 1 && len(current.TopMoods) == 0 {
		topMoodShift = fmt.Sprintf("%s -> %s", previous.TopMoods[0].MoodTag, "not enough data")
	} else if len(previous.TopMoods) == 0 && len(current.TopMoods) >= 1 {
		topMoodShift = fmt.Sprintf("%s -> %s", "not enough data", current.TopMoods[0].MoodTag)
	} else if len(previous.TopMoods) >= 1 && len(current.TopMoods) >= 1 {
		topMoodShift = fmt.Sprintf("%s -> %s", previous.TopMoods[0].MoodTag, current.TopMoods[0].MoodTag)
	} else {
		topMoodShift = "not enough data -> not enough data"
	}

	var topMoodDelta float64
	if len(previous.TopMoods) != 0 {
		var previousMood models.MoodTagFrequency
		for _, mood := range previous.TopMoods {
			if strings.EqualFold(mood.MoodTag, current.TopMoods[0].MoodTag) {
				previousMood = mood
				break
			}
		}
		topMoodDelta = ((current.TopMoods[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
	}

	var topPositiveMoodChange string
	if len(previous.TopMoodsPositiveDays) != 0 {
		var previousMood models.MoodTagFrequency
		for _, mood := range previous.TopMoodsPositiveDays {
			if strings.EqualFold(mood.MoodTag, current.TopMoodsPositiveDays[0].MoodTag) {
				previousMood = mood
				break
			}
		}
		percentChange := ((current.TopMoodsPositiveDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topPositiveMoodChange = fmt.Sprintf("%s %f", current.TopMoodsPositiveDays[0].MoodTag, percentChange)
	}
	var topNegativeMoodChange string
	if len(previous.TopMoodsNegativeDays) != 0 {
		var previousMood models.MoodTagFrequency
		for _, mood := range previous.TopMoodsNegativeDays {
			if strings.EqualFold(mood.MoodTag, current.TopMoodsNegativeDays[0].MoodTag) {
				previousMood = mood
				break
			}
		}
		percentChange := ((current.TopMoodsNegativeDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topNegativeMoodChange = fmt.Sprintf("%s %f", current.TopMoodsNegativeDays[0].MoodTag, percentChange)
	}

	positiveDaysDelta := len(current.PositiveDays) - len(previous.PositiveDays)

	negativeDaysDelta := len(current.NegativeDays) - len(previous.NegativeDays)

	currentTotalEntries := servcie.journalRepository.JournalEntries(current.UserID, current.StartDate, current.EndDate)
	previousTotalEntries := servcie.journalRepository.JournalEntries(previous.UserID, previous.StartDate, previous.EndDate)

	var currentPositiveRatio float64
	if len(current.PositiveDays) != 0 {
		currentPositiveRatio = float64(len(currentTotalEntries)) / float64(len(current.PositiveDays))
	}

	var previoiusPositiveRatio float64
	if len(previous.PositiveDays) != 0 {
		previoiusPositiveRatio = float64(len(previousTotalEntries)) / float64(len(previous.PositiveDays))
	}

	positiveRatioChange := float64(currentPositiveRatio - previoiusPositiveRatio)

	diffs := models.Diff{
		AvgMoodChange:             avgMoodChange,
		TrendChange:               trendChange,
		StabilityChange:           stabilityChange,
		VolatilityDelta:           volatilityDelta,
		VolatilityDeltaPercentage: volatilityDeltaPercentage,
		TopMoodShift:              topMoodShift,
		TopMoodDelta:              topMoodDelta,
		TopPositiveMoodChange:     topPositiveMoodChange,
		TopNegativeMoodChange:     topNegativeMoodChange,
		PositiveDaysDelta:         positiveDaysDelta,
		NegativeDaysDelta:         negativeDaysDelta,
		PositiveRatioChange:       positiveRatioChange,
	}

	return diffs

}
