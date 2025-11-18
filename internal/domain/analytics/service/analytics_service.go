package service

import (
	"fmt"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/utils"

	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/model"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/repository"
)

type AnalyticsService struct {
	analyticsRepository repository.IAnalyticsRepository
}

func NewAnalyticsService(analyticsRepository repository.IAnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepository: analyticsRepository,
	}
}

func (as *AnalyticsService) AnalyzeMood(userID int, startDate time.Time, endDate time.Time) (*model.MoodMetric, error) {

	numDays := utils.NumDaysBetween(startDate, endDate)

	var numDaysPreceding string

	const (
		maxWeekly  = 7
		maxMonthly = 30
		max3Months = 90
	)

	switch {
	case numDays <= maxWeekly:
		numDaysPreceding = "3"
	case numDays <= maxMonthly:
		numDaysPreceding = "7"
	case numDays <= max3Months:
		numDaysPreceding = "14"
	default:
		numDaysPreceding = "3"
	}

	movingAverages, err := as.analyticsRepository.MovingAverages(userID, startDate, endDate, numDaysPreceding)
	if err != nil {
		return nil, fmt.Errorf("failed to get moving averages: %w", err)
	}
	fmt.Println(movingAverages)

	var movingAvg float64
	if len(movingAverages) > 0 {
		lastIndex := len(movingAverages) - 1
		movingAvg = movingAverages[lastIndex].MovingAvg
	}

	moodTrend := utils.Trend(movingAverages)

	standardDeviation, err := as.analyticsRepository.StandardDeviationMood(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get standard deviation: %w", err)
	}

	const (
		noData          = 0
		minModerateMood = 1.5
		minVolatileMood = 3
	)

	stability := utils.StdDeviation(standardDeviation, noData, minModerateMood, minVolatileMood)

	avgMoodRating, err := as.analyticsRepository.AvgMoodRating(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get average mood rating: %w", err)
	}

	mtfPeriod, err := as.analyticsRepository.MoodTagStats(userID, startDate, endDate)
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
	positiveDays, err := as.analyticsRepository.Days(userID, startDate, endDate, greaterThanOrEual, minMoodRatingPositiveDay, positiveMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get positive days: %w", err)
	}
	mtfPositiveDays := utils.MoodTagFrequencies(positiveDays)

	// Neutral Days
	neutralDays, err := as.analyticsRepository.Days(userID, startDate, endDate, equalTo, neutralDayMoodRating, neutralMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get neutral days: %w", err)
	}
	mtfNeutralDays := utils.MoodTagFrequencies(neutralDays)

	// Negative Days
	negativeDays, err := as.analyticsRepository.Days(userID, startDate, endDate, lessThanOrEqual, maxMoodRatingNegativeDay, negativeMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get negative days: %w", err)
	}
	mtfNegativeDays := utils.MoodTagFrequencies(negativeDays)

	// Clinical Days
	clinicalDays, err := as.analyticsRepository.Days(userID, startDate, endDate, greaterThanOrEual, minClinicalMoodRating, clinicalMoodCategory, tagPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinical days: %w", err)
	}
	mtfClinicalDays := utils.MoodTagFrequencies(clinicalDays)

	granularity := utils.Granularity(numDays)

	moodMetrics := &model.MoodMetric{
		UserID:             userID,
		Granularity:        granularity,
		StartDate:          startDate,
		EndDate:            endDate,
		MovingAvg:          movingAvg,
		Trend:              moodTrend,
		StdDeviation:       standardDeviation,
		Stability:          stability,
		AvgRating:          avgMoodRating,
		TopTagOverall:      utils.TopTagStat(mtfPeriod),
		TopTagPositiveDays: utils.TopTagStat(mtfPositiveDays),
		TopTagNegativeDays: utils.TopTagStat(mtfNegativeDays),
		TopTagNeutralDays:  utils.TopTagStat(mtfNeutralDays),
		TopTagClinicalDays: utils.TopTagStat(mtfClinicalDays),
	}

	return moodMetrics, nil
}

func (as *AnalyticsService) AnalyzeSleep(userID int, startDate time.Time, endDate time.Time) (*model.SleepMetric, error) {
	avgSleepHours, err := as.analyticsRepository.AvgSleepHours(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get average sleep hours: %w", err)
	}

	numDays := utils.NumDaysBetween(startDate, endDate)

	var numDaysPreceding string

	const (
		maxWeekly  = 7
		maxMonthly = 30
		max3Months = 90
	)

	switch {
	case numDays <= maxWeekly:
		numDaysPreceding = "3"
	case numDays <= maxMonthly:
		numDaysPreceding = "7"
	case numDays <= max3Months:
		numDaysPreceding = "14"
	default:
		numDaysPreceding = "3"
	}

	movingAverages, err := as.analyticsRepository.MovingAvgSleep(userID, startDate, endDate, numDaysPreceding)
	if err != nil {
		return nil, fmt.Errorf("failed to get moving average sleep: %w", err)
	}

	var movingAvg float64
	if len(movingAverages) > 0 {
		lastIndex := len(movingAverages) - 1
		movingAvg = movingAverages[lastIndex].MovingAvg
	}

	sleepTrend := utils.Trend(movingAverages)

	standardDeviation, err := as.analyticsRepository.StandardDeviationSleep(userID, startDate, endDate)
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

	topSleepQualityTags, err := as.analyticsRepository.SleepQualityTagStats(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get sleep quality tag stats: %w", err)
	}

	// Get day-of-week patterns
	dayOfWeekPatterns, err := as.analyticsRepository.DayOfWeekSleepPatterns(userID, startDate, endDate)
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

	sleepMetric := &model.SleepMetric{
		UserID:             userID,
		Granularity:        granularity,
		StartDate:          startDate,
		EndDate:            endDate,
		AvgSleepHours:      avgSleepHours,
		MovingAvg:          movingAvg,
		SleepTrend:         sleepTrend,
		StdDeviation:       standardDeviation,
		Stability:          stability,
		BestSleepDay:       bestSleepDay,
		WorstSleepDay:      worstSleepDay,
		TopSleepQualityTag: utils.TopTagStat(topSleepQualityTags),
	}

	return sleepMetric, nil
}

func (as *AnalyticsService) AnalyzeMedication(userID int, startDate time.Time, endDate time.Time) (*model.MedicationMetric, error) {
	numDays := utils.NumDaysBetween(startDate, endDate)
	granularity := utils.Granularity(numDays)

	// Get overview statistics
	adherenceRate, err := as.analyticsRepository.MedicationOverviewStats(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get overview stats: %w", err)
	}

	// Get detailed medication statistics
	medicationStats, err := as.analyticsRepository.MedicationDetailedStats(userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication stats: %w", err)
	}

	medicationMetric := &model.MedicationMetric{
		UserID:          userID,
		Granularity:     granularity,
		StartDate:       startDate,
		EndDate:         endDate,
		AdherenceRate:   adherenceRate,
		MedicationStats: medicationStats,
	}

	return medicationMetric, nil
}
