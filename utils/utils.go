package utils

import (
	"regexp"
	"slices"
	"strings"

	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
	"time"
)

func MatchURL(pattern string, url string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(url)
}

func GetUserIDFromPath(path string) string {
	splitPath := strings.Split(path, "/")
	userIDIndex := slices.Index(splitPath, "users") + 1
	return splitPath[userIDIndex]
}

func MoodTagFrequencies(days []models.Day) []models.TagFrequency {
	var tags []string
	for _, day := range days {
		for _, mtf := range day.MoodTagFrequencies {
			tags = append(tags, mtf.TagName)
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

	var moodTagFrequencies []models.TagFrequency
	for key, val := range freq {
		count := int(val)
		percentage := (val / float64(len(tags))) * 100.0

		mtf := models.TagFrequency{
			Count:      count,
			TagName:    key,
			Percentage: percentage,
		}

		moodTagFrequencies = append(moodTagFrequencies, mtf)
	}

	slices.SortFunc(moodTagFrequencies, func(a, b models.TagFrequency) int {
		return int(b.Percentage - a.Percentage)
	})

	if moodTagFrequencies == nil {
		return make([]models.TagFrequency, 0)
	}

	return moodTagFrequencies
}

func FindPreviousMood(currentMoods, previousMoods []models.TagFrequency) models.TagFrequency {
	var previousMood models.TagFrequency
	for _, mood := range previousMoods {
		if strings.EqualFold(mood.TagName, currentMoods[0].TagName) {
			previousMood = mood
			break
		}
	}
	return previousMood
}

func PreviousDates(startDate string, endDate string) (string, string) {
	const layout = "2006-01-02"
	startDateParsed, _ := time.Parse(layout, startDate)
	endDateParsed, _ := time.Parse(layout, endDate)
	diff := endDateParsed.Sub(startDateParsed)
	numDays := int(diff.Hours() / 24)
	previousStartDate := startDateParsed.AddDate(0, 0, -numDays).Format(layout)
	previousEndDate := startDateParsed.AddDate(0, 0, -1).Format(layout)
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

func NumDaysBetween(startDate string, endDate string) int {
	const layout = "2006-01-02" // Correct Go layout
	startDateParsed, _ := time.Parse(layout, startDate)
	endDateParsed, _ := time.Parse(layout, endDate)
	diff := endDateParsed.Sub(startDateParsed)
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
