package analytics

import (
	"fmt"
	"strconv"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/utils"

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
	if len(movingAverages) > 0 {
		lastIndex := len(movingAverages) - 1
		movingAvg = movingAverages[lastIndex].MovingAvg
	}

	moodTrend := utils.Trend(movingAverages)

	standardDeviation := service.moodLogRepository.StandardDeviation(userID, startDate, endDate)

	const (
		noData          = 0
		minModerateMood = 1.5
		minVolatileMood = 3
	)

	stability := utils.StdDeviation(standardDeviation, noData, minModerateMood, minVolatileMood)

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

	utils.AddSleepLogsToDays(userID, service.sleepLogRepository, positiveDays)

	mtfPositiveDays := utils.MoodTagFrequencies(positiveDays)

	neutralDays := service.moodLogRepository.Days(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)

	utils.AddSleepLogsToDays(userID, service.sleepLogRepository, neutralDays)

	mtfNeutralDays := utils.MoodTagFrequencies(neutralDays)

	negativeDays := service.moodLogRepository.Days(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)

	utils.AddSleepLogsToDays(userID, service.sleepLogRepository, negativeDays)

	mtfNegativeDays := utils.MoodTagFrequencies(negativeDays)

	clinicalDays := service.moodLogRepository.Days(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)

	utils.AddSleepLogsToDays(userID, service.sleepLogRepository, clinicalDays)

	mtfClinicalDays := utils.MoodTagFrequencies(clinicalDays)

	positiveStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, greaterThanOrEual, minMoodRatingPositiveDay, positiveMoodCategory, tagPercentage)

	neutralStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)

	negativeStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)

	clinicalStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)

	granularity := utils.Granularity(numDays)

	moodMetrics := &models.MoodMetric{
		UserID:                   userID,
		Granularity:              granularity,
		StartDate:                startDate,
		EndDate:                  endDate,
		MovingAvg:                movingAvg,
		MoodTrend:                moodTrend,
		StdDeviation:             standardDeviation,
		Stability:                stability,
		AvgMoodRating:            avgMoodRating,
		MoodTagStats:             mtfPeriod,
		MoodTagStatsPositiveDays: mtfPositiveDays,
		MoodTagStatsNeutralDays:  mtfNeutralDays,
		MoodTagStatsNegativeDays: mtfNegativeDays,
		MoodTagStatsClinicalDays: mtfClinicalDays,
		PositiveStreaks:          positiveStreaks,
		NeutralStreaks:           neutralStreaks,
		NegativeStreaks:          negativeStreaks,
		ClinicalStreaks:          clinicalStreaks,
		PositiveDays:             positiveDays,
		NeutralDays:              neutralDays,
		NegativeDays:             negativeDays,
		ClinicalDays:             clinicalDays,
		MoodDiffs:                models.MoodDiff{},
	}

	return moodMetrics
}

func (service *analyticsService) moodDiffs(currentPeriod, previousPeriod *models.MoodMetric) models.MoodDiff {

	var avgMoodPercentChange float64
	if previousPeriod.AvgMoodRating > 0 {
		avgMoodPercentChange = utils.PercentChange(currentPeriod.AvgMoodRating, previousPeriod.AvgMoodRating)
	}

	var trendShift string
	if previousPeriod.MoodTrend != "" {
		trendShift = fmt.Sprintf("%s -> %s", previousPeriod.MoodTrend, currentPeriod.MoodTrend)
	}

	var movingAvgPercentChange float64
	if previousPeriod.MovingAvg > 0 {
		movingAvgPercentChange = utils.PercentChange(currentPeriod.MovingAvg, previousPeriod.MovingAvg)
	}

	var stabilityShift string
	if previousPeriod.Stability != "" {
		stabilityShift = fmt.Sprintf("%s -> %s", previousPeriod.Stability, currentPeriod.Stability)
	}

	var stabilityPercentChange float64
	if previousPeriod.StdDeviation > 0 {
		stabilityPercentChange = utils.PercentChange(currentPeriod.StdDeviation, previousPeriod.StdDeviation)

	}

	const moodTagStatsIndex = 0

	var topMoodShift string
	if utils.BothContainValues(currentPeriod.MoodTagStats, previousPeriod.MoodTagStats) {
		previousMood := previousPeriod.MoodTagStats[moodTagStatsIndex]
		currentMood := currentPeriod.MoodTagStats[moodTagStatsIndex]
		topMoodShift = fmt.Sprintf("%s -> %s", previousMood.TagName, currentMood.TagName)
	}

	var topMoodPercentChange string
	if utils.BothContainValues(currentPeriod.MoodTagStats, previousPeriod.MoodTagStats) {
		previousMood := utils.FindPreviousMood(currentPeriod.MoodTagStats, previousPeriod.MoodTagStats)
		currentMood := currentPeriod.MoodTagStats[moodTagStatsIndex]
		topMoodPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, utils.PercentChange(currentMood.Percentage, previousMood.Percentage))
	}

	var topMoodPositiveDaysPercentChange string
	if utils.BothContainValues(currentPeriod.MoodTagStatsPositiveDays, previousPeriod.MoodTagStatsPositiveDays) {
		previousMood := utils.FindPreviousMood(currentPeriod.MoodTagStatsPositiveDays, previousPeriod.MoodTagStatsPositiveDays)
		currentMood := currentPeriod.MoodTagStatsPositiveDays[moodTagStatsIndex]
		percentChange := utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodPositiveDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodNeutralDaysPercentChange string
	if utils.BothContainValues(currentPeriod.MoodTagStatsNeutralDays, previousPeriod.MoodTagStatsNeutralDays) {
		previousMood := utils.FindPreviousMood(currentPeriod.MoodTagStatsNeutralDays, previousPeriod.MoodTagStatsNeutralDays)
		currentMood := currentPeriod.MoodTagStatsNeutralDays[moodTagStatsIndex]
		percentChange := utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodNeutralDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodNegativeDaysPercentChange string
	if utils.BothContainValues(currentPeriod.MoodTagStatsNegativeDays, previousPeriod.MoodTagStatsNegativeDays) {
		previousMood := utils.FindPreviousMood(currentPeriod.MoodTagStatsNegativeDays, previousPeriod.MoodTagStatsNegativeDays)
		currentMood := currentPeriod.MoodTagStatsNegativeDays[moodTagStatsIndex]
		percentChange := utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodNegativeDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var topMoodClinicalDaysPercentChange string
	if utils.BothContainValues(currentPeriod.MoodTagStatsClinicalDays, previousPeriod.MoodTagStatsClinicalDays) {
		previousMood := utils.FindPreviousMood(currentPeriod.MoodTagStatsClinicalDays, previousPeriod.MoodTagStatsClinicalDays)
		currentMood := currentPeriod.MoodTagStatsClinicalDays[moodTagStatsIndex]
		percentChange := utils.PercentChange(currentMood.Percentage, previousMood.Percentage)
		topMoodClinicalDaysPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, percentChange)
	}

	var positiveDaysChange int
	if utils.BothContainValues(currentPeriod.PositiveDays, previousPeriod.PositiveDays) {
		positiveDaysChange = utils.DifferenceInLength(currentPeriod.PositiveDays, previousPeriod.PositiveDays)
	}

	var neutralDaysChange int
	if utils.BothContainValues(currentPeriod.NeutralDays, previousPeriod.NeutralDays) {
		neutralDaysChange = utils.DifferenceInLength(currentPeriod.NeutralDays, previousPeriod.NeutralDays)
	}

	var negativeDaysChange int
	if utils.BothContainValues(currentPeriod.NegativeDays, previousPeriod.NegativeDays) {
		negativeDaysChange = utils.DifferenceInLength(currentPeriod.NegativeDays, previousPeriod.NegativeDays)
	}

	var clinicalDaysChange int
	if utils.BothContainValues(currentPeriod.ClinicalDays, previousPeriod.ClinicalDays) {
		clinicalDaysChange = utils.DifferenceInLength(currentPeriod.ClinicalDays, previousPeriod.ClinicalDays)
	}

	var positiveStreakChange int
	if utils.BothContainValues(currentPeriod.PositiveStreaks, previousPeriod.PositiveStreaks) {
		positiveStreakChange = utils.DifferenceInLength(currentPeriod.PositiveStreaks, previousPeriod.PositiveStreaks)
	}

	var neutralStreakChange int
	if utils.BothContainValues(currentPeriod.NeutralStreaks, previousPeriod.NeutralStreaks) {
		neutralStreakChange = utils.DifferenceInLength(currentPeriod.NeutralStreaks, previousPeriod.NeutralStreaks)
	}

	var negativeStreakChange int
	if utils.BothContainValues(currentPeriod.NegativeStreaks, previousPeriod.NegativeStreaks) {
		negativeStreakChange = utils.DifferenceInLength(currentPeriod.NegativeStreaks, previousPeriod.NegativeStreaks)
	}

	var clinicalStreakChange int
	if utils.BothContainValues(currentPeriod.ClinicalDays, previousPeriod.ClinicalDays) {
		clinicalStreakChange = utils.DifferenceInLength(currentPeriod.ClinicalStreaks, previousPeriod.ClinicalStreaks)
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

	numDays := utils.NumDaysBetween(startDate, endDate)
	numDaysPreceding := strconv.Itoa(numDays)

	movingAverages := service.sleepLogRepository.MovingAvgSleep(userID, startDate, endDate, numDaysPreceding)

	var movingAvg float64
	if len(movingAverages) > 0 {
		lastIndex := len(movingAverages) - 1
		movingAvg = movingAverages[lastIndex].MovingAvg
	}

	sleepTrend := utils.Trend(movingAverages)

	standardDeviation := service.sleepLogRepository.StandardDeviation(userID, startDate, endDate)

	const (
		noData           = 0
		minModerateSleep = 0.5 // 30 mins
		minVolatileSleep = 1.5 // 90 mins
	)

	stability := utils.StdDeviation(standardDeviation, noData, minModerateSleep, minVolatileSleep)

	granularity := utils.Granularity(numDays)

	topSleepQualityTags := service.sleepLogRepository.SleepQualityTagStat(userID, startDate, endDate)

	sleepMetrics := &models.SleepMetric{
		UserID:               userID,
		Granularity:          granularity,
		StartDate:            startDate,
		EndDate:              endDate,
		AvgSleepHours:        avgSleepHours,
		MovingAvg:            movingAvg,
		SleepTrend:           sleepTrend,
		StdDeviation:         standardDeviation,
		Stability:            stability,
		SleepQualityTagStats: topSleepQualityTags,
	}

	return sleepMetrics

}
