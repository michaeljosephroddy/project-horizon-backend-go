package analytics

import (
	"fmt"
	"slices"
	"strconv"

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

	numDays := utils.NumDaysBetween(startDate, endDate)
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

	// TODO fix magic strings
	positiveDays := service.moodLogRepository.Days(userID, startDate, endDate, ">=", "6", "1", "50")
	mtfPositiveDays := utils.MoodTagFrequencies(positiveDays)

	// TODO fix magic strings
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

	granularity := utils.Granularity(numDays)

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

func (service *analyticsService) moodDiffs(currentPeriod, previousPeriod *models.MoodMetric) models.MoodDiff {

	var avgMoodPercentChange float64
	if previousPeriod.AvgMoodRating != 0.0 {
		avgMoodPercentChange = utils.PercentChange(currentPeriod.AvgMoodRating, previousPeriod.AvgMoodRating)
	}

	trendShift := fmt.Sprintf("%s -> %s", previousPeriod.MoodTrend, currentPeriod.MoodTrend)

	var movingAvgPercentChange float64
	if previousPeriod.MovingAvg != 0.0 {
		movingAvgPercentChange = utils.PercentChange(currentPeriod.MovingAvg, previousPeriod.MovingAvg)
	}

	stabilityShift := fmt.Sprintf("%s -> %s", previousPeriod.Stability, currentPeriod.Stability)

	var stabilityPercentChange float64
	if previousPeriod.StdDeviation != 0.0 {
		stabilityPercentChange = utils.PercentChange(currentPeriod.StdDeviation, previousPeriod.StdDeviation)
	}

	// TODO wirte utility method
	var topMoodShift string
	if len(previousPeriod.TopMoods) >= 1 && len(currentPeriod.TopMoods) == 0 {
		topMoodShift = fmt.Sprintf("%s -> %s", previousPeriod.TopMoods[0].TagName, "not enough data")
	} else if len(previousPeriod.TopMoods) == 0 && len(currentPeriod.TopMoods) >= 1 {
		topMoodShift = fmt.Sprintf("%s -> %s", "not enough data", currentPeriod.TopMoods[0].TagName)
	} else if len(previousPeriod.TopMoods) >= 1 && len(currentPeriod.TopMoods) >= 1 {
		topMoodShift = fmt.Sprintf("%s -> %s", previousPeriod.TopMoods[0].TagName, currentPeriod.TopMoods[0].TagName)
	} else {
		topMoodShift = "not enough data -> not enough data"
	}

	topMoodIndex := 0

	var topMoodPercentChange string
	if utils.BothContainValues(currentPeriod.TopMoods, previousPeriod.TopMoods) {
		previousMood := utils.FindMood(currentPeriod.TopMoods, previousPeriod.TopMoods)
		currentMood := currentPeriod.TopMoods[topMoodIndex]
		topMoodPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, utils.PercentChange(currentMood.Percentage, previousMood.Percentage))
	}

	var topMoodPositiveDaysPercentChange string
	if utils.BothContainValues(currentPeriod.TopMoodsPositiveDays, previousPeriod.TopMoodsPositiveDays) {
		previousMood := utils.FindMood(currentPeriod.TopMoodsPositiveDays, previousPeriod.TopMoodsPositiveDays)
		currentMood := currentPeriod.TopMoodsPositiveDays[topMoodIndex]
		percentChange := utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodPositiveDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodNeutralDaysPercentChange string
	if utils.BothContainValues(currentPeriod.TopMoodsNeutralDays, previousPeriod.TopMoodsNeutralDays) {
		previousMood := utils.FindMood(currentPeriod.TopMoodsNeutralDays, previousPeriod.TopMoodsNeutralDays)
		currentMood := currentPeriod.TopMoodsNeutralDays[topMoodIndex]
		percentChange := utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodNeutralDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodNegativeDaysPercentChange string
	if utils.BothContainValues(currentPeriod.TopMoodsNegativeDays, previousPeriod.TopMoodsNegativeDays) {
		previousMood := utils.FindMood(currentPeriod.TopMoodsNegativeDays, previousPeriod.TopMoodsNegativeDays)
		currentMood := currentPeriod.TopMoodsNegativeDays[topMoodIndex]
		percentChange := utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodNegativeDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodClinicalDaysPercentChange string
	if utils.BothContainValues(currentPeriod.TopMoodsClinicalDays, previousPeriod.TopMoodsClinicalDays) {
		previousMood := utils.FindMood(currentPeriod.TopMoodsClinicalDays, previousPeriod.TopMoodsClinicalDays)
		currentMood := currentPeriod.TopMoodsClinicalDays[topMoodIndex]
		percentChange := utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodClinicalDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	// TODO same here could break out repetetive len() - len()
	positiveDaysChange := utils.DifferenceInLength(currentPeriod.PositiveDays, previousPeriod.PositiveDays)
	neutralDaysChange := utils.DifferenceInLength(currentPeriod.NeutralDays, previousPeriod.NeutralDays)
	negativeDaysChange := utils.DifferenceInLength(currentPeriod.NeutralDays, previousPeriod.NeutralDays)
	clinicalDaysChange := utils.DifferenceInLength(currentPeriod.ClinicalDays, previousPeriod.ClinicalDays)

	// TODO same here could break out repetetive len() - len()
	longestPositiveStreakChange := utils.DifferenceInLength(currentPeriod.PositiveStreaks, previousPeriod.PositiveStreaks)
	longestNeutralStreakChange := utils.DifferenceInLength(currentPeriod.NeutralStreaks, previousPeriod.NeutralStreaks)
	longestNegativeStreakChange := utils.DifferenceInLength(currentPeriod.NegativeStreaks, previousPeriod.NegativeStreaks)
	longestClinicalStreakChange := utils.DifferenceInLength(currentPeriod.ClinicalStreaks, previousPeriod.ClinicalStreaks)

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

	numDays := utils.NumDaysBetween(startDate, endDate)
	numDaysPreceding := strconv.Itoa(numDays)

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

	granularity := utils.Granularity(numDays)

	// topSleepQualityTags := service.sleepLogRepository.SleepQualityTagFrequency(userID, startDate, endDate)

	sleepMetrics := &models.SleepMetric{
		UserID:        userID,
		Granularity:   granularity,
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
