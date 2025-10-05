package analytics

import (
	"fmt"
	"strconv"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/analytics_utils"

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

	numDays := analytics_utils.NumDaysBetween(startDate, endDate)
	numDaysPreceding := strconv.Itoa(numDays)

	movingAverages := service.moodLogRepository.MovingAverages(userID, startDate, endDate, numDaysPreceding)

	var movingAvg float64
	if len(movingAverages) > 0 {
		lastIndex := len(movingAverages) - 1
		movingAvg = movingAverages[lastIndex].MovingAvg
	}

	moodTrend := analytics_utils.Trend(movingAverages)

	standardDeviation := service.moodLogRepository.StandardDeviation(userID, startDate, endDate)

	const (
		noData          = 0
		minModerateMood = 1.5
		minVolatileMood = 3
	)

	stability := analytics_utils.StdDeviation(standardDeviation, noData, minModerateMood, minVolatileMood)

	avgMoodRating := service.moodLogRepository.AvgMoodRating(userID, startDate, endDate)

	mtfPeriod := service.moodLogRepository.MoodTagFrequencies(userID, startDate, endDate)

	const (
		greaterThanOrEual        = ">="
		equalTo                  = "="
		lessThanOrEqual          = "<="
		minMoodRatingPositiveDay = "6"
		minClinicalMoodRating    = "1"
		neutralDayMoodRating     = "5"
		maxMoodRatingNegativeDay = "4"
		negativeMoodCategory     = "2"
		positiveMoodCategory     = "1"
		clinicalMoodCategory     = "5"
		neutralMoodCategory      = "3"
		tagPercentage            = "50"
	)

	positiveDays := service.moodLogRepository.Days(userID, startDate, endDate, greaterThanOrEual, minMoodRatingPositiveDay, positiveMoodCategory, tagPercentage)

	mtfPositiveDays := analytics_utils.MoodTagFrequencies(positiveDays)

	neutralDays := service.moodLogRepository.Days(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)

	mtfNeutralDays := analytics_utils.MoodTagFrequencies(neutralDays)

	negativeDays := service.moodLogRepository.Days(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)

	mtfNegativeDays := analytics_utils.MoodTagFrequencies(negativeDays)

	clinicalDays := service.moodLogRepository.Days(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)

	mtfClinicalDays := analytics_utils.MoodTagFrequencies(clinicalDays)

	positiveStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, greaterThanOrEual, minMoodRatingPositiveDay, positiveMoodCategory, tagPercentage)

	neutralStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)

	negativeStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)

	clinicalStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)

	granularity := analytics_utils.Granularity(numDays)

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
	if previousPeriod.AvgMoodRating > 0 {
		avgMoodPercentChange = analytics_utils.PercentChange(currentPeriod.AvgMoodRating, previousPeriod.AvgMoodRating)
	}

	var trendShift string
	if previousPeriod.MoodTrend != "" {
		trendShift = fmt.Sprintf("%s -> %s", previousPeriod.MoodTrend, currentPeriod.MoodTrend)
	}

	var movingAvgPercentChange float64
	if previousPeriod.MovingAvg > 0 {
		movingAvgPercentChange = analytics_utils.PercentChange(currentPeriod.MovingAvg, previousPeriod.MovingAvg)
	}

	var stabilityShift string
	if previousPeriod.Stability != "" {
		stabilityShift = fmt.Sprintf("%s -> %s", previousPeriod.Stability, currentPeriod.Stability)
	}

	var stabilityPercentChange float64
	if previousPeriod.StdDeviation > 0 {
		stabilityPercentChange = analytics_utils.PercentChange(currentPeriod.StdDeviation, previousPeriod.StdDeviation)

	}

	const topMoodIndex = 0

	var topMoodShift string
	if analytics_utils.BothContainValues(currentPeriod.TopMoods, previousPeriod.TopMoods) {
		previousMood := previousPeriod.TopMoods[topMoodIndex]
		currentMood := currentPeriod.TopMoods[topMoodIndex]
		topMoodShift = fmt.Sprintf("%s -> %s", previousMood.TagName, currentMood.TagName)
	}

	var topMoodPercentChange string
	if analytics_utils.BothContainValues(currentPeriod.TopMoods, previousPeriod.TopMoods) {
		previousMood := analytics_utils.FindPreviousMood(currentPeriod.TopMoods, previousPeriod.TopMoods)
		currentMood := currentPeriod.TopMoods[topMoodIndex]
		topMoodPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, analytics_utils.PercentChange(currentMood.Percentage, previousMood.Percentage))
	}

	var topMoodPositiveDaysPercentChange string
	if analytics_utils.BothContainValues(currentPeriod.TopMoodsPositiveDays, previousPeriod.TopMoodsPositiveDays) {
		previousMood := analytics_utils.FindPreviousMood(currentPeriod.TopMoodsPositiveDays, previousPeriod.TopMoodsPositiveDays)
		currentMood := currentPeriod.TopMoodsPositiveDays[topMoodIndex]
		percentChange := analytics_utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodPositiveDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodNeutralDaysPercentChange string
	if analytics_utils.BothContainValues(currentPeriod.TopMoodsNeutralDays, previousPeriod.TopMoodsNeutralDays) {
		previousMood := analytics_utils.FindPreviousMood(currentPeriod.TopMoodsNeutralDays, previousPeriod.TopMoodsNeutralDays)
		currentMood := currentPeriod.TopMoodsNeutralDays[topMoodIndex]
		percentChange := analytics_utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodNeutralDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodNegativeDaysPercentChange string
	if analytics_utils.BothContainValues(currentPeriod.TopMoodsNegativeDays, previousPeriod.TopMoodsNegativeDays) {
		previousMood := analytics_utils.FindPreviousMood(currentPeriod.TopMoodsNegativeDays, previousPeriod.TopMoodsNegativeDays)
		currentMood := currentPeriod.TopMoodsNegativeDays[topMoodIndex]
		percentChange := analytics_utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodNegativeDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodClinicalDaysPercentChange string
	if analytics_utils.BothContainValues(currentPeriod.TopMoodsClinicalDays, previousPeriod.TopMoodsClinicalDays) {
		previousMood := analytics_utils.FindPreviousMood(currentPeriod.TopMoodsClinicalDays, previousPeriod.TopMoodsClinicalDays)
		currentMood := currentPeriod.TopMoodsClinicalDays[topMoodIndex]
		percentChange := analytics_utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodClinicalDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var positiveDaysChange int
	if analytics_utils.BothContainValues(currentPeriod.PositiveDays, previousPeriod.PositiveDays) {
		positiveDaysChange = analytics_utils.DifferenceInLength(currentPeriod.PositiveDays, previousPeriod.PositiveDays)
	}

	var neutralDaysChange int
	if analytics_utils.BothContainValues(currentPeriod.NeutralDays, previousPeriod.NeutralDays) {
		neutralDaysChange = analytics_utils.DifferenceInLength(currentPeriod.NeutralDays, previousPeriod.NeutralDays)
	}

	var negativeDaysChange int
	if analytics_utils.BothContainValues(currentPeriod.NegativeDays, previousPeriod.NegativeDays) {
		negativeDaysChange = analytics_utils.DifferenceInLength(currentPeriod.NegativeDays, previousPeriod.NegativeDays)
	}

	var clinicalDaysChange int
	if analytics_utils.BothContainValues(currentPeriod.ClinicalDays, previousPeriod.ClinicalDays) {
		clinicalDaysChange = analytics_utils.DifferenceInLength(currentPeriod.ClinicalDays, previousPeriod.ClinicalDays)
	}

	var positiveStreakChange int
	if analytics_utils.BothContainValues(currentPeriod.PositiveStreaks, previousPeriod.PositiveStreaks) {
		positiveStreakChange = analytics_utils.DifferenceInLength(currentPeriod.PositiveStreaks, previousPeriod.PositiveStreaks)
	}

	var neutralStreakChange int
	if analytics_utils.BothContainValues(currentPeriod.NeutralStreaks, previousPeriod.NeutralStreaks) {
		neutralStreakChange = analytics_utils.DifferenceInLength(currentPeriod.NeutralStreaks, previousPeriod.NeutralStreaks)
	}

	var negativeStreakChange int
	if analytics_utils.BothContainValues(currentPeriod.NegativeStreaks, previousPeriod.NegativeStreaks) {
		negativeStreakChange = analytics_utils.DifferenceInLength(currentPeriod.NegativeStreaks, previousPeriod.NegativeStreaks)
	}

	var clinicalStreakChange int
	if analytics_utils.BothContainValues(currentPeriod.ClinicalDays, previousPeriod.ClinicalDays) {
		clinicalStreakChange = analytics_utils.DifferenceInLength(currentPeriod.ClinicalStreaks, previousPeriod.ClinicalStreaks)
	}

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
		PositiveStreakChange:             positiveStreakChange,
		NeutralStreakChange:              neutralStreakChange,
		NegativeStreakChange:             negativeStreakChange,
		ClinicalStreakChange:             clinicalStreakChange,
	}

	return moodDiffs
}

func (service *analyticsService) analyzeSleep(userID string, startDate string, endDate string) *models.SleepMetric {

	avgSleepHours := service.sleepLogRepository.AvgSleepHours(userID, startDate, endDate)

	numDays := analytics_utils.NumDaysBetween(startDate, endDate)
	numDaysPreceding := strconv.Itoa(numDays)

	movingAverages := service.sleepLogRepository.MovingAvgSleep(userID, startDate, endDate, numDaysPreceding)

	var movingAvg float64
	if len(movingAverages) > 0 {
		lastIndex := len(movingAverages) - 1
		movingAvg = movingAverages[lastIndex].MovingAvg
	}

	sleepTrend := analytics_utils.Trend(movingAverages)

	standardDeviation := service.sleepLogRepository.StandardDeviation(userID, startDate, endDate)

	const (
		noData           = 0
		minModerateSleep = 0.5 // 30 mins
		minVolatileSleep = 1.5 // 90 mins
	)

	stability := analytics_utils.StdDeviation(standardDeviation, noData, minModerateSleep, minVolatileSleep)

	granularity := analytics_utils.Granularity(numDays)

	topSleepQualityTags := service.sleepLogRepository.SleepQualityTagFrequency(userID, startDate, endDate)

	sleepMetrics := &models.SleepMetric{
		UserID:              userID,
		Granularity:         granularity,
		StartDate:           startDate,
		EndDate:             endDate,
		AvgSleepHours:       avgSleepHours,
		MovingAvg:           movingAvg,
		SleepTrend:          sleepTrend,
		StdDeviation:        standardDeviation,
		Stability:           stability,
		TopSleepQualityTags: topSleepQualityTags,
	}

	return sleepMetrics

}
