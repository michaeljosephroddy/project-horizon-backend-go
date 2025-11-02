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

func (service *analyticsService) analyzeMood(userID string, startDate time.Time, endDate time.Time) (*models.MoodMetric, error) {

	numDays := utils.NumDaysBetween(startDate, endDate)
	numDaysPreceding := strconv.Itoa(numDays)

	movingAverages, err := service.moodLogRepository.MovingAverages(userID, startDate, endDate, numDaysPreceding)
	if err != nil {
		return nil, fmt.Errorf("failed to get moving averages: %w", err)
	}

	var movingAvg float64
	if len(movingAverages) > 0 {
		lastIndex := len(movingAverages) - 1
		movingAvg = movingAverages[lastIndex].MovingAvg
	}

	moodTrend := utils.Trend(movingAverages)

	standardDeviation, err := service.moodLogRepository.StandardDeviation(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get standard deviation: %w", err)
	}

	const (
		noData          = 0
		minModerateMood = 1.5
		minVolatileMood = 3
	)

	stability := utils.StdDeviation(standardDeviation, noData, minModerateMood, minVolatileMood)

	avgMoodRating, err := service.moodLogRepository.AvgMoodRating(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get average mood rating: %w", err)
	}

	mtfPeriod, err := service.moodLogRepository.MoodTagFrequencies(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get mood tag frequencies: %w", err)
	}

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
	positiveDays, err := service.moodLogRepository.Days(userID, startDate, endDate, greaterThanOrEual, minMoodRatingPositiveDay, positiveMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get positive days: %w", err)
	}
	if err := utils.AddSleepLogsToDays(userID, service.sleepLogRepository, positiveDays); err != nil {
		return nil, fmt.Errorf("failed to add sleep logs to positive days: %w", err)
	}
	if err := utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, positiveDays); err != nil {
		return nil, fmt.Errorf("failed to add medication logs to positive days: %w", err)
	}
	mtfPositiveDays := utils.MoodTagFrequencies(positiveDays)

	// Neutral Days
	neutralDays, err := service.moodLogRepository.Days(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get neutral days: %w", err)
	}
	if err := utils.AddSleepLogsToDays(userID, service.sleepLogRepository, neutralDays); err != nil {
		return nil, fmt.Errorf("failed to add sleep logs to neutral days: %w", err)
	}
	if err := utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, neutralDays); err != nil {
		return nil, fmt.Errorf("failed to add medication logs to neutral days: %w", err)
	}
	mtfNeutralDays := utils.MoodTagFrequencies(neutralDays)

	// Negative Days
	negativeDays, err := service.moodLogRepository.Days(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get negative days: %w", err)
	}
	if err := utils.AddSleepLogsToDays(userID, service.sleepLogRepository, negativeDays); err != nil {
		return nil, fmt.Errorf("failed to add sleep logs to negative days: %w", err)
	}
	if err := utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, negativeDays); err != nil {
		return nil, fmt.Errorf("failed to add medication logs to negative days: %w", err)
	}
	mtfNegativeDays := utils.MoodTagFrequencies(negativeDays)

	// Clinical Days
	clinicalDays, err := service.moodLogRepository.Days(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinical days: %w", err)
	}
	if err := utils.AddSleepLogsToDays(userID, service.sleepLogRepository, clinicalDays); err != nil {
		return nil, fmt.Errorf("failed to add sleep logs to clinical days: %w", err)
	}
	if err := utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, clinicalDays); err != nil {
		return nil, fmt.Errorf("failed to add medication logs to clinical days: %w", err)
	}
	mtfClinicalDays := utils.MoodTagFrequencies(clinicalDays)

	// Positive Streaks
	positiveStreaks, err := service.moodLogRepository.Streaks(userID, startDate, endDate, greaterThanOrEual, minMoodRatingPositiveDay, positiveMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get positive streaks: %w", err)
	}
	for _, streak := range positiveStreaks {
		if err := utils.AddSleepLogsToDays(userID, service.sleepLogRepository, streak.Days); err != nil {
			return nil, fmt.Errorf("failed to add sleep logs to positive streak: %w", err)
		}
		if err := utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, streak.Days); err != nil {
			return nil, fmt.Errorf("failed to add medication logs to positive streak: %w", err)
		}
	}

	// Neutral Streaks
	neutralStreaks, err := service.moodLogRepository.Streaks(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get neutral streaks: %w", err)
	}
	for _, streak := range neutralStreaks {
		if err := utils.AddSleepLogsToDays(userID, service.sleepLogRepository, streak.Days); err != nil {
			return nil, fmt.Errorf("failed to add sleep logs to neutral streak: %w", err)
		}
		if err := utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, streak.Days); err != nil {
			return nil, fmt.Errorf("failed to add medication logs to neutral streak: %w", err)
		}
	}

	// Negative Streaks
	negativeStreaks, err := service.moodLogRepository.Streaks(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get negative streaks: %w", err)
	}
	for _, streak := range negativeStreaks {
		if err := utils.AddSleepLogsToDays(userID, service.sleepLogRepository, streak.Days); err != nil {
			return nil, fmt.Errorf("failed to add sleep logs to negative streak: %w", err)
		}
		if err := utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, streak.Days); err != nil {
			return nil, fmt.Errorf("failed to add medication logs to negative streak: %w", err)
		}
	}

	// Clinical Streaks
	clinicalStreaks, err := service.moodLogRepository.Streaks(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinical streaks: %w", err)
	}
	for _, streak := range clinicalStreaks {
		if err := utils.AddSleepLogsToDays(userID, service.sleepLogRepository, streak.Days); err != nil {
			return nil, fmt.Errorf("failed to add sleep logs to clinical streak: %w", err)
		}
		if err := utils.AddMedicationLogsToDays(userID, service.medicationLogRepository, streak.Days); err != nil {
			return nil, fmt.Errorf("failed to add medication logs to clinical streak: %w", err)
		}
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

	return moodMetrics, nil
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

func (service *analyticsService) analyzeSleep(userID string, startDate time.Time, endDate time.Time) (*models.SleepMetric, error) {
	avgSleepHours, err := service.sleepLogRepository.AvgSleepHours(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get average sleep hours: %w", err)
	}

	numDays := utils.NumDaysBetween(startDate, endDate)
	numDaysPreceding := strconv.Itoa(numDays)
	movingAverages, err := service.sleepLogRepository.MovingAvgSleep(userID, startDate, endDate, numDaysPreceding)
	if err != nil {
		return nil, fmt.Errorf("failed to get moving average sleep: %w", err)
	}

	var movingAvg float64
	if len(movingAverages) > 0 {
		lastIndex := len(movingAverages) - 1
		movingAvg = movingAverages[lastIndex].MovingAvg
	}

	sleepTrend := utils.Trend(movingAverages)

	standardDeviation, err := service.sleepLogRepository.StandardDeviation(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get sleep standard deviation: %w", err)
	}

	const (
		noData           = 0
		minModerateSleep = 0.5 // 30 mins
		minVolatileSleep = 1.5 // 90 mins
	)
	stability := utils.StdDeviation(standardDeviation, noData, minModerateSleep, minVolatileSleep)

	granularity := utils.Granularity(numDays)

	topSleepQualityTags, err := service.sleepLogRepository.SleepQualityTagStat(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get sleep quality tag stats: %w", err)
	}

	// Get day-of-week patterns
	dayOfWeekPatterns, err := service.sleepLogRepository.DayOfWeekSleepPatterns(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get day-of-week sleep patterns: %w", err)
	}

	// Calculate best/worst days by duration only
	var bestSleepDay, worstSleepDay string

	if len(dayOfWeekPatterns) > 0 {
		maxHours := -1.0
		minHours := 999.0

		for _, pattern := range dayOfWeekPatterns {
			if pattern.AvgSleepHours > maxHours {
				maxHours = pattern.AvgSleepHours
				bestSleepDay = pattern.DayOfWeek
			}
			if pattern.AvgSleepHours < minHours {
				minHours = pattern.AvgSleepHours
				worstSleepDay = pattern.DayOfWeek
			}
		}
	}

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
		BestSleepDay:         bestSleepDay,
		WorstSleepDay:        worstSleepDay,
		SleepDiffs:           models.SleepDiff{},
	}

	return sleepMetrics, nil
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
