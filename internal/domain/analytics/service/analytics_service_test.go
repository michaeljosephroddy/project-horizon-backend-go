package service

import (
	"testing"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAnalyticsRepository is a mock implementation of repository.AnalyticsRepositoryInterface
type MockAnalyticsRepository struct {
	mock.Mock
}

// Ensure MockAnalyticsRepository implements the interface
var _ interface {
	MovingAverages(userID int, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error)
	StandardDeviation(userID int, startDate time.Time, endDate time.Time) (float64, error)
	AvgMoodRating(userID int, startDate time.Time, endDate time.Time) (float64, error)
	MoodTagFrequencies(userID int, startDate time.Time, endDate time.Time) ([]model.TagStat, error)
	Days(userID int, startDate time.Time, endDate time.Time, operator string, threshold string, category string, tagPercentage string) ([]model.Day, error)
	AvgSleepHours(userID int, startDate time.Time, endDate time.Time) (float64, error)
	MovingAvgSleep(userID int, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error)
	SleepStandardDeviation(userID int, startDate time.Time, endDate time.Time) (float64, error)
	SleepQualityTagStat(userID int, startDate time.Time, endDate time.Time) ([]model.TagStat, error)
	DayOfWeekSleepPatterns(userID int, startDate time.Time, endDate time.Time) ([]model.DayOfWeekSleepPattern, error)
	OverviewStats(userID int, startDate time.Time, endDate time.Time) (float64, error)
	MedicationDetailedStats(userID int, startDate time.Time, endDate time.Time) ([]model.MedicationStats, error)
} = (*MockAnalyticsRepository)(nil)

func (m *MockAnalyticsRepository) MovingAverages(userID int, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error) {
	args := m.Called(userID, startDate, endDate, numDaysPreceding)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.MovingAverage), args.Error(1)
}

func (m *MockAnalyticsRepository) StandardDeviation(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockAnalyticsRepository) AvgMoodRating(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockAnalyticsRepository) MoodTagFrequencies(userID int, startDate time.Time, endDate time.Time) ([]model.TagStat, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TagStat), args.Error(1)
}

func (m *MockAnalyticsRepository) Days(userID int, startDate time.Time, endDate time.Time, operator string, threshold string, category string, tagPercentage string) ([]model.Day, error) {
	args := m.Called(userID, startDate, endDate, operator, threshold, category, tagPercentage)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Day), args.Error(1)
}

func (m *MockAnalyticsRepository) AvgSleepHours(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockAnalyticsRepository) MovingAvgSleep(userID int, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error) {
	args := m.Called(userID, startDate, endDate, numDaysPreceding)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.MovingAverage), args.Error(1)
}

func (m *MockAnalyticsRepository) SleepStandardDeviation(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockAnalyticsRepository) SleepQualityTagStat(userID int, startDate time.Time, endDate time.Time) ([]model.TagStat, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TagStat), args.Error(1)
}

func (m *MockAnalyticsRepository) DayOfWeekSleepPatterns(userID int, startDate time.Time, endDate time.Time) ([]model.DayOfWeekSleepPattern, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.DayOfWeekSleepPattern), args.Error(1)
}

func (m *MockAnalyticsRepository) OverviewStats(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockAnalyticsRepository) MedicationDetailedStats(userID int, startDate time.Time, endDate time.Time) ([]model.MedicationStats, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.MedicationStats), args.Error(1)
}

// Test AnalyzeMood
func TestAnalyzeMood_Success(t *testing.T) {
	mockRepo := new(MockAnalyticsRepository)
	service := NewAnalyticsService(mockRepo)

	userID := 123
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)

	movingAvgs := []model.MovingAverage{
		{Date: startDate, MovingAvg: 5.0},
		{Date: startDate.AddDate(0, 0, 1), MovingAvg: 5.5},
	}

	positiveDays := []model.Day{
		{Date: startDate, DailyAvgRating: 7.0},
	}

	tagStats := []model.TagStat{
		{TagName: "happy", Count: 10, Percentage: 50.0},
	}

	mockRepo.On("MovingAverages", userID, startDate, endDate, "6").Return(movingAvgs, nil)
	mockRepo.On("StandardDeviation", userID, startDate, endDate).Return(1.2, nil)
	mockRepo.On("AvgMoodRating", userID, startDate, endDate).Return(6.5, nil)
	mockRepo.On("MoodTagFrequencies", userID, startDate, endDate).Return(tagStats, nil)
	mockRepo.On("Days", userID, startDate, endDate, ">=", "6", "1", "50").Return(positiveDays, nil)
	mockRepo.On("Days", userID, startDate, endDate, "=", "5", "3", "50").Return([]model.Day{}, nil)
	mockRepo.On("Days", userID, startDate, endDate, "<=", "4", "2", "50").Return([]model.Day{}, nil)
	mockRepo.On("Days", userID, startDate, endDate, ">=", "1", "5", "50").Return([]model.Day{}, nil)

	result, err := service.AnalyzeMood(userID, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, 5.5, result.MovingAvg)
	assert.Equal(t, 6.5, result.AvgRating)
	assert.Equal(t, 1.2, result.StdDeviation)
	mockRepo.AssertExpectations(t)
}

// Test AnalyzeSleep
func TestAnalyzeSleep_Success(t *testing.T) {
	mockRepo := new(MockAnalyticsRepository)
	service := NewAnalyticsService(mockRepo)

	userID := 123
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)

	movingAvgs := []model.MovingAverage{
		{Date: startDate, MovingAvg: 7.5},
		{Date: startDate.AddDate(0, 0, 1), MovingAvg: 8.0},
	}

	dayPatterns := []model.DayOfWeekSleepPattern{
		{DayOfWeek: "Monday", DayNumber: 2, AvgSleepHours: 7.0, TotalEntries: 4},
		{DayOfWeek: "Tuesday", DayNumber: 3, AvgSleepHours: 8.5, TotalEntries: 4},
		{DayOfWeek: "Wednesday", DayNumber: 4, AvgSleepHours: 6.5, TotalEntries: 4},
	}

	tagStats := []model.TagStat{
		{TagName: "good", Count: 15, Percentage: 60.0},
	}

	mockRepo.On("AvgSleepHours", userID, startDate, endDate).Return(7.5, nil)
	mockRepo.On("MovingAvgSleep", userID, startDate, endDate, "6").Return(movingAvgs, nil)
	mockRepo.On("SleepStandardDeviation", userID, startDate, endDate).Return(0.8, nil)
	mockRepo.On("SleepQualityTagStat", userID, startDate, endDate).Return(tagStats, nil)
	mockRepo.On("DayOfWeekSleepPatterns", userID, startDate, endDate).Return(dayPatterns, nil)

	result, err := service.AnalyzeSleep(userID, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, 7.5, result.AvgSleepHours)
	assert.Equal(t, 8.0, result.MovingAvg)
	assert.Equal(t, 0.8, result.StdDeviation)
	assert.Equal(t, "Tuesday", result.BestSleepDay)
	assert.Equal(t, "Wednesday", result.WorstSleepDay)
	mockRepo.AssertExpectations(t)
}

// Test AnalyzeMedication
func TestAnalyzeMedication_Success(t *testing.T) {
	mockRepo := new(MockAnalyticsRepository)
	service := NewAnalyticsService(mockRepo)

	userID := 123
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)

	medicationStats := []model.MedicationStats{
		{
			MedicationID:        1,
			Name:                "Aspirin",
			TotalDoses:          7,
			DaysActive:          7,
			AvgDosesPerDay:      1.0,
			AvgTakenAtTime:      "08:00:00",
			TimingStdDevMinutes: 15.0,
			TimingDescription:   "8:00 AM ± 15 minutes",
			EarliestTime:        "07:45:00",
			LatestTime:          "08:15:00",
			LongestStreak:       7,
			CurrentStreak:       7,
		},
	}

	mockRepo.On("OverviewStats", userID, startDate, endDate).Return(85.5, nil)
	mockRepo.On("MedicationDetailedStats", userID, startDate, endDate).Return(medicationStats, nil)

	result, err := service.AnalyzeMedication(userID, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, 85.5, result.AdherenceRate)
	assert.Equal(t, 1, len(result.MedicationStats))
	assert.Equal(t, "Aspirin", result.MedicationStats[0].Name)
	mockRepo.AssertExpectations(t)
}
