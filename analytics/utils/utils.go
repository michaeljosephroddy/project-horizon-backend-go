package utils

import (
	"fmt"
	"slices"
	"strings"

	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
	"github.com/michaeljosephroddy/project-horizon-backend-go/database"
)

// TODO make these package private if not used outside analytics package
func MoodTagFrequencies(days []models.Day) []models.TagStat {
	var tags []string
	for _, day := range days {
		for _, mtf := range day.MoodLogs {
			tags = append(tags, mtf.MoodTags...)
		}
	}

	const (
		zeroVal      = 0.0
		incrementVal = 1.0
	)

	freq := make(map[string]float64)
	for _, tag := range tags {
		if _, exists := freq[tag]; !exists {
			freq[tag] = zeroVal
		}

		freq[tag] = freq[tag] + incrementVal
	}

	var moodTagFrequencies []models.TagStat
	for key, val := range freq {
		count := int(val)
		percentage := (val / float64(len(tags))) * 100.0

		mtf := models.TagStat{
			Count:      count,
			TagName:    key,
			Percentage: percentage,
		}

		moodTagFrequencies = append(moodTagFrequencies, mtf)
	}

	slices.SortFunc(moodTagFrequencies, func(a, b models.TagStat) int {
		return int(b.Percentage - a.Percentage)
	})

	if moodTagFrequencies == nil {
		return make([]models.TagStat, 0)
	}

	return moodTagFrequencies
}

func FindPreviousMood(currentMoods, previousMoods []models.TagStat) models.TagStat {
	var previousMood models.TagStat
	for _, mood := range previousMoods {
		if strings.EqualFold(mood.TagName, currentMoods[0].TagName) {
			previousMood = mood
			break
		}
	}
	return previousMood
}

func PreviousDates(startDate time.Time, endDate time.Time) (time.Time, time.Time) {
	const layout = "2006-01-02"
	diff := endDate.Sub(startDate)
	numDays := int(diff.Hours() / 24)
	previousStartDate, _ := time.Parse(layout, startDate.AddDate(0, 0, -numDays).Format(layout))
	previousEndDate, _ := time.Parse(layout, startDate.AddDate(0, 0, -1).Format(layout))
	return previousStartDate, previousEndDate
}

func Trend(movingAvergaes []models.MovingAverage) string {
	const (
		increasing = "increasing"
		decreasing = "decreasing"
		flat       = "flat"
	)

	var trend string
	if len(movingAvergaes) >= 2 {
		lastIndex := len(movingAvergaes) - 1
		secondLastIndex := len(movingAvergaes) - 2

		last := movingAvergaes[lastIndex]
		prev := movingAvergaes[secondLastIndex]

		switch {
		case last.MovingAvg > prev.MovingAvg:
			trend = increasing
		case last.MovingAvg < prev.MovingAvg:
			trend = decreasing
		default:
			trend = flat
		}

	} else {
		trend = ""
	}

	return trend
}

func Granularity(numDays int) string {
	const (
		maxWeekly   = 7
		maxMonthly  = 28
		max3Months  = 84
		weekly      = "weekly"
		monthly     = "monthly"
		threeMonths = "3-months"
		custom      = "custom"
	)

	var granularity string

	switch {
	case numDays <= maxWeekly:
		granularity = weekly
	case numDays <= maxMonthly:
		granularity = monthly
	case numDays <= max3Months:
		granularity = threeMonths
	default:
		granularity = custom
	}

	return granularity
}

func NumDaysBetween(startDate time.Time, endDate time.Time) int {
	diff := endDate.Sub(startDate)
	return int(diff.Hours() / 24)
}

func PercentChange(a, b float64) float64 {
	return ((a - b) / b) * 100
}

func BothContainValues[T any](a, b []T) bool {
	return len(a) != 0 && len(b) != 0
}

func DifferenceInLength[T any](a, b []T) int {
	return len(a) - len(b)
}

func StdDeviation(standardDeviation float64, noData float64, minModerateVal float64, minVolatileVal float64) string {
	const (
		stable   = "stable"
		moderate = "moderate"
		volatile = "volatile"
	)

	var stability string

	switch {
	case standardDeviation == noData:
		stability = "" // e.g., only 1 data point
	case standardDeviation < minModerateVal:
		stability = stable
	case standardDeviation < minVolatileVal:
		stability = moderate
	default:
		stability = volatile

	}
	return stability
}

func AddSleepLogsToDays(userID string, slr *database.SleepLogRepository, days []models.Day) error {
	for i, day := range days {
		sleepLogs, err := slr.SleepLogs(userID, day.Date, day.Date)
		if err != nil {
			return fmt.Errorf("failed to get sleep logs for day %s: %w", day.Date.Format("2006-01-02"), err)
		}
		if len(sleepLogs) == 0 {
			days[i].SleepLogs = make([]models.SleepLog, 0)
			continue
		}
		// Just assign all the logs directly since the query already filtered by date
		days[i].SleepLogs = sleepLogs
	}
	return nil
}

func AddMedicationLogsToDays(userID string, mlr *database.MedicationLogRepository, days []models.Day) error {
	for i, day := range days {
		medicationLogs, err := mlr.MedicationLogs(userID, day.Date, day.Date)
		if err != nil {
			return fmt.Errorf("failed to get medication logs for day %s: %w", day.Date.Format("2006-01-02"), err)
		}
		if len(medicationLogs) == 0 {
			days[i].MedicationLogs = make([]models.MedicationLog, 0)
			continue
		}
		// Just assign all the logs directly since the query already filtered by date
		days[i].MedicationLogs = medicationLogs
	}
	return nil
}
