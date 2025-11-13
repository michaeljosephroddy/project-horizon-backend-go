package repository

import (
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/model"
	"time"
)

// AnalyticsRepositoryInterface defines the contract for analytics repository operations
type AnalyticsRepositoryInterface interface {
	MovingAverages(userID string, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error)
	StandardDeviation(userID string, startDate time.Time, endDate time.Time) (float64, error)
	AvgMoodRating(userID string, startDate time.Time, endDate time.Time) (float64, error)
	MoodTagFrequencies(userID string, startDate time.Time, endDate time.Time) ([]model.TagStat, error)
	Days(userID string, startDate time.Time, endDate time.Time, operator string, threshold string, category string, tagPercentage string) ([]model.Day, error)
	AvgSleepHours(userID string, startDate time.Time, endDate time.Time) (float64, error)
	MovingAvgSleep(userID string, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error)
	SleepStandardDeviation(userID string, startDate time.Time, endDate time.Time) (float64, error)
	SleepQualityTagStat(userID string, startDate time.Time, endDate time.Time) ([]model.TagStat, error)
	DayOfWeekSleepPatterns(userID string, startDate time.Time, endDate time.Time) ([]model.DayOfWeekSleepPattern, error)
	OverviewStats(userID string, startDate time.Time, endDate time.Time) (float64, error)
	MedicationDetailedStats(userID string, startDate time.Time, endDate time.Time) ([]model.MedicationStats, error)
}
