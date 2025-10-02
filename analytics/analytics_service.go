package analytics

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/utils"

	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type analyticsService struct {
	moodLogRepository  *database.MoodLogRepository
	sleepLogRepository *database.SleepLogRepository
}

func NewAnalyticsService(moodLogRepository *database.MoodLogRepository, sleepLogRepository *database.SleepLogRepository) *analyticsService {
	return &analyticsService{
		moodLogRepository:  moodLogRepository,
		sleepLogRepository: sleepLogRepository,
	}
}

func (service *analyticsService) analyzeMood(userID string, startDate string, endDate string) *models.MoodMetric {

	layout := "2006-01-02" // Correct Go layout
	startDateParsed, _ := time.Parse(layout, startDate)
	endDateParsed, _ := time.Parse(layout, endDate)

	diff := endDateParsed.Sub(startDateParsed)
	numDays := int(diff.Hours() / 24)
	numDaysPreceding := strconv.Itoa(numDays)

	movingAverages := service.moodLogRepository.MovingAverages(userID, startDate, endDate, numDaysPreceding)

	var movingAvg float64
	if len(movingAverages) >= 2 {
		movingAvg = movingAverages[len(movingAverages)-1].MovingAvg
	}

	moodTrend := utils.DetermineTrend(movingAverages)

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
		stability = "volatile"
	}

	avgMoodRating := service.moodLogRepository.AvgMoodRating(userID, startDate, endDate)

	mtfPeriod := service.moodLogRepository.MoodTagFrequencies(userID, startDate, endDate)

	slices.SortFunc(mtfPeriod, func(a, b models.TagFrequency) int {
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

	moodMetrics := &models.MoodMetric{
		UserID:               userID,
		Granularity:          granularity,
		StartDate:            startDate,
		EndDate:              endDate,
		MovingAvg:            movingAvg,
		MoodTrend:            moodTrend,
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
		MoodDiffs:            models.MoodDiff{},
	}

	return moodMetrics
}

func (service *analyticsService) moodDiffs(current, previous *models.MoodMetric) models.MoodDiff {

	var avgMoodPercentChange float64
	if previous.AvgMoodRating != 0.0 {
		avgMoodPercentChange = ((current.AvgMoodRating - previous.AvgMoodRating) / previous.AvgMoodRating) * 100
	}

	trendShift := fmt.Sprintf("%s -> %s", previous.MoodTrend, current.MoodTrend)

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
		topMoodShift = fmt.Sprintf("%s -> %s", previous.TopMoods[0].TagName, "not enough data")
	} else if len(previous.TopMoods) == 0 && len(current.TopMoods) >= 1 {
		topMoodShift = fmt.Sprintf("%s -> %s", "not enough data", current.TopMoods[0].TagName)
	} else if len(previous.TopMoods) >= 1 && len(current.TopMoods) >= 1 {
		topMoodShift = fmt.Sprintf("%s -> %s", previous.TopMoods[0].TagName, current.TopMoods[0].TagName)
	} else {
		topMoodShift = "not enough data -> not enough data"
	}

	var topMoodPercentChange string
	if len(previous.TopMoods) != 0 && len(current.TopMoods) != 0 {
		previousMood := utils.FindMood(current.TopMoods, previous.TopMoods)
		topMoodPercentChange = fmt.Sprintf("%s %f", current.TopMoods[0].TagName, ((current.TopMoods[0].Percentage-previousMood.Percentage)/previousMood.Percentage)*100)
	}

	var topMoodPositiveDaysPercentChange string
	if len(previous.TopMoodsPositiveDays) != 0 && len(current.TopMoodsPositiveDays) != 0 {
		previousMood := utils.FindMood(current.TopMoodsPositiveDays, previous.TopMoodsPositiveDays)
		percentChange := ((current.TopMoodsPositiveDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topMoodPositiveDaysPercentChange = fmt.Sprintf("%s %f", current.TopMoodsPositiveDays[0].TagName, percentChange)
	}

	var topMoodNeutralDaysPercentChange string
	if len(previous.TopMoodsNeutralDays) != 0 && len(current.TopMoodsNeutralDays) != 0 {
		previousMood := utils.FindMood(current.TopMoodsNeutralDays, previous.TopMoodsNeutralDays)
		percentChange := ((current.TopMoodsNeutralDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topMoodNeutralDaysPercentChange = fmt.Sprintf("%s %f", current.TopMoodsNeutralDays[0].TagName, percentChange)
	}

	var topMoodNegativeDaysPercentChange string
	if len(previous.TopMoodsNegativeDays) != 0 && len(current.TopMoodsNegativeDays) != 0 {
		previousMood := utils.FindMood(current.TopMoodsNegativeDays, previous.TopMoodsNegativeDays)
		percentChange := ((current.TopMoodsNegativeDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topMoodNegativeDaysPercentChange = fmt.Sprintf("%s %f", current.TopMoodsNegativeDays[0].TagName, percentChange)
	}

	var topMoodClinicalDaysPercentChange string
	if len(previous.TopMoodsClinicalDays) != 0 && len(current.TopMoodsClinicalDays) != 0 {
		previousMood := utils.FindMood(current.TopMoodsClinicalDays, previous.TopMoodsClinicalDays)
		percentChange := ((current.TopMoodsClinicalDays[0].Percentage - previousMood.Percentage) / previousMood.Percentage) * 100
		topMoodClinicalDaysPercentChange = fmt.Sprintf("%s %f", current.TopMoodsClinicalDays[0].TagName, percentChange)
	}

	positiveDaysChange := len(current.PositiveDays) - len(previous.PositiveDays)
	neutralDaysChange := len(current.NeutralDays) - len(previous.NeutralDays)
	negativeDaysChange := len(current.NegativeDays) - len(previous.NegativeDays)
	clinicalDaysChange := len(current.ClinicalDays) - len(previous.ClinicalDays)

	longestPositiveStreakChange := len(current.PositiveStreaks) - len(previous.PositiveStreaks)
	longestNeutralStreakChange := len(current.NeutralStreaks) - len(previous.NeutralStreaks)
	longestNegativeStreakChange := len(current.NegativeStreaks) - len(previous.NegativeStreaks)
	longestClinicalStreakChange := len(current.ClinicalStreaks) - len(previous.ClinicalStreaks)

	moodDiffs := models.MoodDiff{
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

	return moodDiffs
}

func (service *analyticsService) analyzeSleep(userID string, startDate string, endDate string) *models.SleepMetric {

	avgSleepHours := service.sleepLogRepository.AvgSleepHours(userID, startDate, endDate)

	layout := "2006-01-02"
	startParsed, _ := time.Parse(layout, startDate)
	endParsed, _ := time.Parse(layout, endDate)
	duration := endParsed.Sub(startParsed)
	numDaysPreceding := strconv.Itoa(int(duration.Hours() / 24))

	movingAverages := service.sleepLogRepository.MovingAvgSleep(userID, startDate, endDate, numDaysPreceding)

	var movingAvg float64
	if len(movingAverages) != 1 {
		movingAvg = movingAverages[len(movingAverages)-1].MovingAvg
	}

	sleepTrend := utils.DetermineTrend(movingAverages)

	standardDeviation := service.sleepLogRepository.StandardDeviation(userID, startDate, endDate)

	var stability string

	switch {
	case standardDeviation == 0:
		stability = "not enough data" // e.g., only 1 data point
	case standardDeviation < 0.5: //30 mins
		stability = "stable"
	case standardDeviation < 1.5: //90 mins
		stability = "moderate"
	default:
		stability = "volatile"
	}

	// topSleepQualityTags := service.sleepLogRepository.SleepQualityTagFrequency(userID, startDate, endDate)

	sleepMetrics := &models.SleepMetric{
		UserID:        userID,
		StartDate:     startDate,
		EndDate:       endDate,
		AvgSleepHours: avgSleepHours,
		MovingAvg:     movingAvg,
		SleepTrend:    sleepTrend,
		StdDeviation:  standardDeviation,
		Stability:     stability,
	}

	return sleepMetrics

}
