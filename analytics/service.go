package analytics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/utils"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
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

	var avgRating models.MetricChange
	if previousPeriod.AvgRating > 0 {
		avgRating.PercentChange = utils.PercentChange(currentPeriod.AvgRating, previousPeriod.AvgRating)
	}

	var trend models.ShiftChange
	if previousPeriod.Trend != "" {
		trend.Description = fmt.Sprintf("%s -> %s", previousPeriod.Trend, currentPeriod.Trend)
		trend.Change = utils.PercentChange(currentPeriod.MovingAvg, previousPeriod.MovingAvg)
	}

	var stability models.MetricChange
	if previousPeriod.Stability != "" {
		stability.Shift = fmt.Sprintf("%s -> %s", previousPeriod.Stability, currentPeriod.Stability)
	}
	if previousPeriod.StdDeviation > 0 {
		stability.PercentChange = utils.PercentChange(currentPeriod.StdDeviation, previousPeriod.StdDeviation)
	}

	const tagStatsIndex = 0

	var topTag models.MetricChange
	if utils.BothContainValues(currentPeriod.TagStats, previousPeriod.TagStats) {
		previousTag := previousPeriod.TagStats[tagStatsIndex]
		currentTag := currentPeriod.TagStats[tagStatsIndex]
		topTag.Shift = fmt.Sprintf("%s -> %s", previousTag.TagName, currentTag.TagName)

		previousTagForCalc := utils.FindPreviousMood(currentPeriod.TagStats, previousPeriod.TagStats)
		topTag.PercentChange = utils.PercentChange(currentTag.Percentage, previousTagForCalc.Percentage)
	}

	// Helper function to calculate category diffs
	calculateCategoryDiff := func(current, previous models.CategoryData) models.CategoryDiff {
		var diff models.CategoryDiff

		if utils.BothContainValues(current.TagStats, previous.TagStats) {
			previousTag := utils.FindPreviousMood(current.TagStats, previous.TagStats)
			currentTag := current.TagStats[tagStatsIndex]
			diff.TopTag.PercentChange = utils.PercentChange(currentTag.Percentage, previousTag.Percentage)

			// Add shift for category top tag
			previousTopTag := previous.TagStats[tagStatsIndex]
			currentTopTag := current.TagStats[tagStatsIndex]
			diff.TopTag.Shift = fmt.Sprintf("%s -> %s", previousTopTag.TagName, currentTopTag.TagName)
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
		AvgRating: avgRating,
		Trend:     trend,
		Stability: stability,
		TopTag:    topTag,
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
	var avgSleepHours models.MetricChange
	if previousPeriod.AvgSleepHours > 0 {
		avgSleepHours.PercentChange = utils.PercentChange(currentPeriod.AvgSleepHours, previousPeriod.AvgSleepHours)
	}

	var trend models.ShiftChange
	if previousPeriod.SleepTrend != "" {
		trend.Description = fmt.Sprintf("%s -> %s", previousPeriod.SleepTrend, currentPeriod.SleepTrend)
		trend.Change = utils.PercentChange(currentPeriod.MovingAvg, previousPeriod.MovingAvg)
	}

	var stability models.MetricChange
	if previousPeriod.Stability != "" {
		stability.Shift = fmt.Sprintf("%s -> %s", previousPeriod.Stability, currentPeriod.Stability)
	}
	if previousPeriod.StdDeviation > 0 {
		stability.PercentChange = utils.PercentChange(currentPeriod.StdDeviation, previousPeriod.StdDeviation)
	}

	const moodTagStatsIndex = 0
	var topQualityTag models.MetricChange
	if utils.BothContainValues(currentPeriod.SleepQualityTagStats, previousPeriod.SleepQualityTagStats) {
		previousMood := previousPeriod.SleepQualityTagStats[moodTagStatsIndex]
		currentMood := currentPeriod.SleepQualityTagStats[moodTagStatsIndex]
		topQualityTag.Shift = fmt.Sprintf("%s -> %s", previousMood.TagName, currentMood.TagName)

		previousMoodForCalc := utils.FindPreviousMood(currentPeriod.SleepQualityTagStats, previousPeriod.SleepQualityTagStats)
		topQualityTag.PercentChange = utils.PercentChange(currentMood.Percentage, previousMoodForCalc.Percentage)
	}

	sleepDiffs := models.SleepDiff{
		AvgSleepHours: avgSleepHours,
		Trend:         trend,
		Stability:     stability,
		TopQualityTag: topQualityTag,
	}

	return sleepDiffs
}
