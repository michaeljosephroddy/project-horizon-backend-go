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

func MoodTagFrequencies(data []models.Day) []models.TagFrequency {
	var tags []string
	for i := 0; i < len(data); i++ {
		for _, val := range data[i].MoodTagFrequencies {
			tags = append(tags, val.TagName)
		}
	}

	freq := make(map[string]float64)
	for _, tag := range tags {
		if _, exists := freq[tag]; !exists {
			freq[tag] = 0.0
		}
		freq[tag] = freq[tag] + 1.0
	}

	var moodTagFrequencies []models.TagFrequency
	for key, val := range freq {
		mtf := models.TagFrequency{
			Count:      int(val),
			TagName:    key,
			Percentage: (val / float64(len(tags))) * 100.0,
		}
		moodTagFrequencies = append(moodTagFrequencies, mtf)
	}

	slices.SortFunc(moodTagFrequencies, func(a, b models.TagFrequency) int {
		if a.Percentage > b.Percentage {
			return -1
		} else if a.Percentage < b.Percentage {
			return 1
		} else {
			return 0
		}
	})

	if moodTagFrequencies == nil {
		return make([]models.TagFrequency, 0)
	}

	return moodTagFrequencies
}

func FindMood(currentMoods, previousMoods []models.TagFrequency) models.TagFrequency {
	var previousMood models.TagFrequency
	for _, mood := range previousMoods {
		if strings.EqualFold(mood.TagName, currentMoods[0].TagName) {
			previousMood = mood
			break
		}
	}

	return previousMood
}

func CalculatePreviousDates(startDate string, endDate string) (string, string) {
	layout := "2006-01-02"
	startDateParsed, _ := time.Parse(layout, startDate)
	endDateParsed, _ := time.Parse(layout, endDate)
	diff := endDateParsed.Sub(startDateParsed)
	numDays := int(diff.Hours() / 24)
	previousStart := startDateParsed.AddDate(0, 0, -numDays).Format(layout)
	previousEnd := startDateParsed.AddDate(0, 0, -1).Format(layout)

	return previousStart, previousEnd
}

func DetermineTrend(data []models.MovingAverage) string {

	var trend string
	if len(data) >= 2 {
		last := data[len(data)-1]
		prev := data[len(data)-2]

		switch {
		case last.MovingAvg > prev.MovingAvg:
			trend = "increasing"
		case last.MovingAvg < prev.MovingAvg:
			trend = "decreasing"
		default:
			trend = "flat"
		}
	} else {
		trend = "not enough data"
	}

	return trend
}
