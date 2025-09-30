package analytics

import (
	"fmt"
	"slices"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/utils"

	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type analyticsService struct {
	moodLogRepository *database.MoodLogRepository
}

func NewAnalyticsService(moodLogRepository *database.MoodLogRepository) *analyticsService {
	return &analyticsService{
		moodLogRepository: moodLogRepository,
	}
}

func (service *analyticsService) metrics(userID string, startDate string, endDate string) *models.Metrics {

	current := service.analyze(userID, startDate, endDate)
	fmt.Println(startDate, " ", endDate, current.MovingAvg)

	previousStart, previousEnd := utils.CalculatePreviousDates(startDate, endDate)

	previous := service.analyze(userID, previousStart, previousEnd)
	fmt.Println(previousStart, " ", previousEnd, previous.MovingAvg)

	diffs := service.diffs(current, previous)

	current.Diffs = diffs

	return current
}

func (service *analyticsService) analyze(userID string, startDate string, endDate string) *models.Metrics {

	movingAverages := service.moodLogRepository.MovingAverages(userID, startDate, endDate)

	var movingAvg float64
	var trend string
	if len(movingAverages) >= 2 {
		last := movingAverages[len(movingAverages)-1]
		prev := movingAverages[len(movingAverages)-2]
		movingAvg = last.MovingAvg

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

	standardDeviation := service.moodLogRepository.StandardDeviation(userID, startDate, endDate)

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

	avgMoodRating := service.moodLogRepository.AvgMoodRating(userID, startDate, endDate)

	mtfPeriod := service.moodLogRepository.MoodTagFrequencies(userID, startDate, endDate)

	slices.SortFunc(mtfPeriod, func(a, b models.MoodTagFrequency) int {
		if a.Percentage > b.Percentage {
			return -1
		} else if a.Percentage < b.Percentage {
			return 1
		} else {
			return 0
		}
	})

	positiveDays := service.moodLogRepository.Days(userID, startDate, endDate, ">=", "6", "1", "50")
	mtfPositiveDays := utils.MoodTagFrequencies(positiveDays)

	neutralDays := service.moodLogRepository.Days(userID, startDate, endDate, "=", "5", "3", "50")
	mtfNeutralDays := utils.MoodTagFrequencies(neutralDays)

	negativeDays := service.moodLogRepository.Days(userID, startDate, endDate, "<=", "4", "2", "50")
	mtfNegativeDays := utils.MoodTagFrequencies(negativeDays)

	clinicalDays := service.moodLogRepository.Days(userID, startDate, endDate, ">=", "1", "5", "50")
	mtfClinicalDays := utils.MoodTagFrequencies(clinicalDays)

	positiveStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, ">=", "6", "1", "50")

	neutralStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, "=", "5", "3", "50")

	negativeStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, "<=", "4", "2", "50")

	clinicalStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, ">=", "1", "5", "50")

	layout := "2006-01-02" // Correct Go layout
	startDateParsed, _ := time.Parse(layout, startDate)
	endDateParsed, _ := time.Parse(layout, endDate)

	diff := endDateParsed.Sub(startDateParsed)
	numDays := int(diff.Hours() / 24)

	var granularity string
	switch {
	case numDays <= 7:
		granularity = "weekly"
	case numDays <= 28:
		granularity = "monthly"
	case numDays <= 84:
		granularity = "3-months"
	default:
		granularity = "custom"
	}

	metrics := &models.Metrics{
		UserID:               userID,
		Granularity:          granularity,
		StartDate:            startDate,
		EndDate:              endDate,
		MovingAvg:            movingAvg,
		Trend:                trend,
		StdDeviation:         standardDeviation,
		Stability:            stability,
		AvgMoodRating:        avgMoodRating,
		TopMoods:             mtfPeriod,
		TopMoodsPositiveDays: mtfPositiveDays,
		TopMoodsNeutralDays:  mtfNeutralDays,
		TopMoodsNegativeDays: mtfNegativeDays,
		TopMoodsClinicalDays: mtfClinicalDays,
		PositiveStreaks:      positiveStreaks,
		NeutralStreaks:       neutralStreaks,
		NegativeStreaks:      negativeStreaks,
		ClinicalStreaks:      clinicalStreaks,
		PositiveDays:         positiveDays,
		NeutralDays:          neutralDays,
		NegativeDays:         negativeDays,
		ClinicalDays:         clinicalDays,
		Diffs:                models.Diff{},
	}

	return metrics
}

func (service *analyticsService) diffs(current, previous *models.Metrics) models.Diff {

	var avgMoodPercentChange float64
	if previous.AvgMoodRating != 0.0 {
		avgMoodPercentChange = ((current.AvgMoodRating - previous.AvgMoodRating) / previous.AvgMoodRating) * 100
	}

	trendShift := fmt.Sprintf("%s -> %s", previous.Trend, current.Trend)

	var movingAvgPercentChange float64
	if previous.MovingAvg != 0.0 {
		fmt.Println(previous.MovingAvg)
		movingAvgPercentChange = ((current.MovingAvg - previous.MovingAvg) / previous.MovingAvg) * 100
	}

	stabilityShift := fmt.Sprintf("%s -> %s", previous.Stability, current.Stability)

	var stabilityPercentChange float64
	if previous.StdDeviation != 0.0 {
		stabilityPercentChange = ((current.StdDeviation - previous.StdDeviation) / previous.StdDeviation) * 100
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

	// TODO fix bug here when user id is 3, problem with accessing index in slice, list is empty so it throws error
	// line 204 complains but will need to check other ifs below because they have same logic pretty much

	var topMoodPercentChange string
	if len(previous.TopMoods) != 0 && len(current.TopMoods) != 0 {
		previousMood := utils.FindMood(current.TopMoods, previous.TopMoods)
		topMoodPercentChange = fmt.Sprintf("%s %f", current.TopMoods[0].MoodTag, ((current.TopMoods[0].Percentage-previousMood.Percentage)/previousMood.Percentage)*100)
	}

	var topMoodPositiveDaysPercentChange string
	if len(previous.TopMoodsPositiveDays) != 0 && len(current.TopMoodsPositiveDays) != 0 {
		previousMood := utils.FindMood(current.TopMoodsPositiveDays, previous.TopMoodsPositiveDays)
		percentChange := ((current.TopMoodsPositiveDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topMoodPositiveDaysPercentChange = fmt.Sprintf("%s %f", current.TopMoodsPositiveDays[0].MoodTag, percentChange)
	}

	var topMoodNeutralDaysPercentChange string
	if len(previous.TopMoodsNeutralDays) != 0 && len(current.TopMoodsNeutralDays) != 0 {
		previousMood := utils.FindMood(current.TopMoodsNeutralDays, previous.TopMoodsNeutralDays)
		percentChange := ((current.TopMoodsNeutralDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topMoodNeutralDaysPercentChange = fmt.Sprintf("%s %f", current.TopMoodsNeutralDays[0].MoodTag, percentChange)
	}

	var topMoodNegativeDaysPercentChange string
	if len(previous.TopMoodsNegativeDays) != 0 && len(current.TopMoodsNegativeDays) != 0 {
		previousMood := utils.FindMood(current.TopMoodsNegativeDays, previous.TopMoodsNegativeDays)
		percentChange := ((current.TopMoodsNegativeDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topMoodNegativeDaysPercentChange = fmt.Sprintf("%s %f", current.TopMoodsNegativeDays[0].MoodTag, percentChange)
	}

	var topMoodClinicalDaysPercentChange string
	if len(previous.TopMoodsClinicalDays) != 0 && len(current.TopMoodsClinicalDays) != 0 {
		previousMood := utils.FindMood(current.TopMoodsClinicalDays, previous.TopMoodsClinicalDays)
		percentChange := ((current.TopMoodsClinicalDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topMoodClinicalDaysPercentChange = fmt.Sprintf("%s %f", current.TopMoodsClinicalDays[0].MoodTag, percentChange)
	}

	positiveDaysChange := len(current.PositiveDays) - len(previous.PositiveDays)
	neutralDaysChange := len(current.NeutralDays) - len(previous.NeutralDays)
	negativeDaysChange := len(current.NegativeDays) - len(previous.NegativeDays)
	clinicalDaysChange := len(current.ClinicalDays) - len(previous.ClinicalDays)

	longestPositiveStreakChange := len(current.PositiveStreaks) - len(previous.PositiveStreaks)
	longestNeutralStreakChange := len(current.NeutralStreaks) - len(previous.NeutralStreaks)
	longestNegativeStreakChange := len(current.NegativeStreaks) - len(previous.NegativeStreaks)
	longestClinicalStreakChange := len(current.ClinicalStreaks) - len(previous.ClinicalStreaks)

	diffs := models.Diff{
		AvgMoodPercentChange:             avgMoodPercentChange,
		TrendShift:                       trendShift,
		MovingAvgPercentChange:           movingAvgPercentChange,
		StabilityShift:                   stabilityShift,
		StabilityPercentChange:           stabilityPercentChange,
		TopMoodShift:                     topMoodShift,
		TopMoodPercentChange:             topMoodPercentChange,
		TopMoodPositiveDaysPercentChange: topMoodPositiveDaysPercentChange,
		TopMoodNeutralDaysPercentChange:  topMoodNeutralDaysPercentChange,
		TopMoodNegativeDaysPercentChange: topMoodNegativeDaysPercentChange,
		TopMoodClinicalDaysPercentChange: topMoodClinicalDaysPercentChange,
		PositiveDaysChange:               positiveDaysChange,
		NeutralDaysChange:                neutralDaysChange,
		NegativeDaysChange:               negativeDaysChange,
		ClinicalDaysChange:               clinicalDaysChange,
		LongestPositiveStreakChange:      longestPositiveStreakChange,
		LongestNeutralStreakChange:       longestNeutralStreakChange,
		LongestNegativeStreakChange:      longestNegativeStreakChange,
		LongestClinicalStreakChange:      longestClinicalStreakChange,
	}

	return diffs

}
