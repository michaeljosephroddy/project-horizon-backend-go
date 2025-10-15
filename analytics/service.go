package analytics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/utils"

	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
)

type analyticsService struct {
	moodLogRepository       *database.MoodLogRepository
	sleepLogRepository      *database.SleepLogRepository
	medicationLogRepository *database.MedicationLogRepository
}

func NewAnalyticsService(moodLogRepository *database.MoodLogRepository, sleepLogRepository *database.SleepLogRepository, medicationLogRepository *database.MedicationLogRepository) *analyticsService {
	return &analyticsService{
		moodLogRepository:       moodLogRepository,
		sleepLogRepository:      sleepLogRepository,
		medicationLogRepository: medicationLogRepository,
	}
}

func (service *analyticsService) analyzeMood(userID string, startDate time.Time, endDate time.Time) *models.MoodMetric {

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

	// Positive Days
	positiveDays := service.moodLogRepository.Days(userID, startDate, endDate, greaterThanOrEual, minMoodRatingPositiveDay, positiveMoodCategory, tagPercentage)
	utils.AddSleepLogsToDays(userID, service.sleepLogRepository, positiveDays)
	utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, positiveDays)
	mtfPositiveDays := utils.MoodTagFrequencies(positiveDays)

	// Neutral Days
	neutralDays := service.moodLogRepository.Days(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)
	utils.AddSleepLogsToDays(userID, service.sleepLogRepository, neutralDays)
	utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, neutralDays)
	mtfNeutralDays := utils.MoodTagFrequencies(neutralDays)

	// Negative Days
	negativeDays := service.moodLogRepository.Days(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)
	utils.AddSleepLogsToDays(userID, service.sleepLogRepository, negativeDays)
	utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, negativeDays)
	mtfNegativeDays := utils.MoodTagFrequencies(negativeDays)

	// Clinical Days
	clinicalDays := service.moodLogRepository.Days(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)
	utils.AddSleepLogsToDays(userID, service.sleepLogRepository, clinicalDays)
	utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, clinicalDays)
	mtfClinicalDays := utils.MoodTagFrequencies(clinicalDays)

	// Positive Streaks
	positiveStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, greaterThanOrEual, minMoodRatingPositiveDay, positiveMoodCategory, tagPercentage)
	for _, streak := range positiveStreaks {
		utils.AddSleepLogsToDays(userID, service.sleepLogRepository, streak.Days)
		utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, streak.Days)
	}

	// Neutral Streaks
	neutralStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)
	for _, streak := range neutralStreaks {
		utils.AddSleepLogsToDays(userID, service.sleepLogRepository, streak.Days)
		utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, streak.Days)
	}

	// Negative Streaks
	negativeStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)
	for _, streak := range negativeStreaks {
		utils.AddSleepLogsToDays(userID, service.sleepLogRepository, streak.Days)
		utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, streak.Days)
	}

	// Clinical Streaks
	clinicalStreaks := service.moodLogRepository.Streaks(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)
	for _, streak := range clinicalStreaks {
		utils.AddSleepLogsToDays(userID, service.sleepLogRepository, streak.Days)
		utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, streak.Days)
	}

	granularity := utils.Granularity(numDays)

	moodMetrics := &models.MoodMetric{
		UserID:       userID,
		Granularity:  granularity,
		StartDate:    startDate,
		EndDate:      endDate,
		MovingAvg:    movingAvg,
		Trend:        moodTrend,
		StdDeviation: standardDeviation,
		Stability:    stability,
		AvgRating:    avgMoodRating,
		TagStats:     mtfPeriod,
		Categories: models.MoodCategories{
			Positive: models.CategoryData{
				TagStats: mtfPositiveDays,
				Streaks:  positiveStreaks,
				Days:     positiveDays,
			},
			Neutral: models.CategoryData{
				TagStats: mtfNeutralDays,
				Streaks:  neutralStreaks,
				Days:     neutralDays,
			},
			Negative: models.CategoryData{
				TagStats: mtfNegativeDays,
				Streaks:  negativeStreaks,
				Days:     negativeDays,
			},
			Clinical: models.CategoryData{
				TagStats: mtfClinicalDays,
				Streaks:  clinicalStreaks,
				Days:     clinicalDays,
			},
		},
		Diffs: models.MoodDiff{},
	}

	return moodMetrics
}

func (service *analyticsService) moodDiffs(currentPeriod, previousPeriod *models.MoodMetric) models.MoodDiff {

	var avgRatingChange float64
	if previousPeriod.AvgRating > 0 {
		avgRatingChange = utils.PercentChange(currentPeriod.AvgRating, previousPeriod.AvgRating)
	}

	var trendShift string
	if previousPeriod.Trend != "" {
		trendShift = fmt.Sprintf("%s -> %s", previousPeriod.Trend, currentPeriod.Trend)
	}

	var movingAvgChange float64
	if previousPeriod.MovingAvg > 0 {
		movingAvgChange = utils.PercentChange(currentPeriod.MovingAvg, previousPeriod.MovingAvg)
	}

	var stabilityShift string
	if previousPeriod.Stability != "" {
		stabilityShift = fmt.Sprintf("%s -> %s", previousPeriod.Stability, currentPeriod.Stability)
	}

	var stabilityChange float64
	if previousPeriod.StdDeviation > 0 {
		stabilityChange = utils.PercentChange(currentPeriod.StdDeviation, previousPeriod.StdDeviation)
	}

	const tagStatsIndex = 0

	var topTagShift string
	if utils.BothContainValues(currentPeriod.TagStats, previousPeriod.TagStats) {
		previousTag := previousPeriod.TagStats[tagStatsIndex]
		currentTag := currentPeriod.TagStats[tagStatsIndex]
		topTagShift = fmt.Sprintf("%s -> %s", previousTag.TagName, currentTag.TagName)
	}

	var topTagChange float64
	if utils.BothContainValues(currentPeriod.TagStats, previousPeriod.TagStats) {
		previousTag := utils.FindPreviousMood(currentPeriod.TagStats, previousPeriod.TagStats)
		currentTag := currentPeriod.TagStats[tagStatsIndex]
		topTagChange = utils.PercentChange(currentTag.Percentage, previousTag.Percentage)
	}

	// Helper function to calculate category diffs
	calculateCategoryDiff := func(current, previous models.CategoryData) models.CategoryDiff {
		var diff models.CategoryDiff

		if utils.BothContainValues(current.TagStats, previous.TagStats) {
			previousTag := utils.FindPreviousMood(current.TagStats, previous.TagStats)
			currentTag := current.TagStats[tagStatsIndex]
			diff.TopTagChange = utils.PercentChange(currentTag.Percentage, previousTag.Percentage)
		}

		if utils.BothContainValues(current.Days, previous.Days) {
			diff.DaysChange = utils.DifferenceInLength(current.Days, previous.Days)
		}

		if utils.BothContainValues(current.Streaks, previous.Streaks) {
			diff.StreakChange = utils.DifferenceInLength(current.Streaks, previous.Streaks)
		}

		return diff
	}

	moodDiffs := models.MoodDiff{
		AvgRatingChange:    avgRatingChange,
		TrendShift:         trendShift,
		MovingAvgChange:    movingAvgChange,
		StabilityShift:     stabilityShift,
		StabilityChange:    stabilityChange,
		TopTagShift:        topTagShift,
		TopTagChange:       topTagChange,
		Categories: models.CategoryDiffs{
			Positive: calculateCategoryDiff(currentPeriod.Categories.Positive, previousPeriod.Categories.Positive),
			Neutral:  calculateCategoryDiff(currentPeriod.Categories.Neutral, previousPeriod.Categories.Neutral),
			Negative: calculateCategoryDiff(currentPeriod.Categories.Negative, previousPeriod.Categories.Negative),
			Clinical: calculateCategoryDiff(currentPeriod.Categories.Clinical, previousPeriod.Categories.Clinical),
		},
	}

	return moodDiffs
}

func (service *analyticsService) analyzeSleep(userID string, startDate time.Time, endDate time.Time) *models.SleepMetric {

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
		SleepDiffs:           models.SleepDiff{},
	}

	return sleepMetrics
}

func (service *analyticsService) sleepDiffs(currentPeriod, previousPeriod *models.SleepMetric) models.SleepDiff {

	var avgSleepPercentChange float64
	if previousPeriod.AvgSleepHours > 0 {
		avgSleepPercentChange = utils.PercentChange(currentPeriod.AvgSleepHours, previousPeriod.AvgSleepHours)
	}

	var trendShift string
	if previousPeriod.SleepTrend != "" {
		trendShift = fmt.Sprintf("%s -> %s", previousPeriod.SleepTrend, currentPeriod.SleepTrend)
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

	var topSleepQualityTagShift string
	if utils.BothContainValues(currentPeriod.SleepQualityTagStats, previousPeriod.SleepQualityTagStats) {
		previousMood := previousPeriod.SleepQualityTagStats[moodTagStatsIndex]
		currentMood := currentPeriod.SleepQualityTagStats[moodTagStatsIndex]
		topSleepQualityTagShift = fmt.Sprintf("%s -> %s", previousMood.TagName, currentMood.TagName)
	}

	var topSleepQualityTagPercentChange string
	if utils.BothContainValues(currentPeriod.SleepQualityTagStats, previousPeriod.SleepQualityTagStats) {
		previousMood := utils.FindPreviousMood(currentPeriod.SleepQualityTagStats, previousPeriod.SleepQualityTagStats)
		currentMood := currentPeriod.SleepQualityTagStats[moodTagStatsIndex]
		topSleepQualityTagPercentChange = fmt.Sprintf("%s %f", currentMood.TagName, utils.PercentChange(currentMood.Percentage, previousMood.Percentage))
	}

	sleepDiffs := models.SleepDiff{
		AvgSleepHoursPercentChange:      avgSleepPercentChange,
		TrendShift:                      trendShift,
		MovingAvgPercentChange:          movingAvgPercentChange,
		StabilityShift:                  stabilityShift,
		StabilityPercentChange:          stabilityPercentChange,
		TopSleepQualityTagShift:         topSleepQualityTagShift,
		TopSleepQualityTagPercentChange: topSleepQualityTagPercentChange,
	}

	return sleepDiffs

}
