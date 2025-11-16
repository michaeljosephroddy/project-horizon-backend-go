package repository

import (
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/model"
	"time"
)

// AnalyticsRepositoryInterface defines the contract for analytics repository operations
type IAnalyticsRepository interface {
	MovingAverages(userID int, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error)
	StandardDeviationMood(userID int, startDate time.Time, endDate time.Time) (float64, error)
	AvgMoodRating(userID int, startDate time.Time, endDate time.Time) (float64, error)
	MoodTagStats(userID int, startDate time.Time, endDate time.Time) ([]model.TagStat, error)
	Days(userID int, startDate time.Time, endDate time.Time, operator string, threshold string, category string, tagPercentage string) ([]model.Day, error)
	AvgSleepHours(userID int, startDate time.Time, endDate time.Time) (float64, error)
	MovingAvgSleep(userID int, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error)
	StandardDeviationSleep(userID int, startDate time.Time, endDate time.Time) (float64, error)
	SleepQualityTagStats(userID int, startDate time.Time, endDate time.Time) ([]model.TagStat, error)
	DayOfWeekSleepPatterns(userID int, startDate time.Time, endDate time.Time) ([]model.DayOfWeekSleepPattern, error)
	MedicationOverviewStats(userID int, startDate time.Time, endDate time.Time) (float64, error)
	MedicationDetailedStats(userID int, startDate time.Time, endDate time.Time) ([]model.MedicationStats, error)
}
