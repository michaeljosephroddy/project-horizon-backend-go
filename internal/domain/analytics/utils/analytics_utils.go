package utils

import (
	"slices"

	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/model"
)

func TopTagStat(data []model.TagStat) model.TagStat {
	if len(data) == 0 {
		return model.TagStat{}
	}
	return data[0]

}

func MoodTagFrequencies(days []model.Day) []model.TagStat {
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

	var moodTagFrequencies []model.TagStat
	for key, val := range freq {
		count := int(val)
		percentage := (val / float64(len(tags))) * 100.0

		mtf := model.TagStat{
			Count:      count,
			TagName:    key,
			Percentage: percentage,
		}

		moodTagFrequencies = append(moodTagFrequencies, mtf)
	}

	slices.SortFunc(moodTagFrequencies, func(a, b model.TagStat) int {
		return int(b.Percentage - a.Percentage)
	})

	if moodTagFrequencies == nil {
		return make([]model.TagStat, 0)
	}

	return moodTagFrequencies
}

func Trend(movingAvergaes []model.MovingAverage) string {
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
		trend = "not enough data"
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
	if b == 0 {
		return 0
	}
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
