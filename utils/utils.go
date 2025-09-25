package utils

import (
	"regexp"
	"slices"
	"strings"

	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
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

func MoodTagFrequencies(data []models.Day) []models.MoodTagFrequency {
	var tags []string
	for i := 0; i < len(data); i++ {
		for _, val := range data[i].MoodTagFrequencies {
			tags = append(tags, val.MoodTag)
		}
	}

	freq := make(map[string]float64)
	for _, tag := range tags {
		if _, exists := freq[tag]; !exists {
			freq[tag] = 0.0
		}
		freq[tag] = freq[tag] + 1.0
	}

	var moodTagFrequencies []models.MoodTagFrequency
	for key, val := range freq {
		mtf := models.MoodTagFrequency{
			Count:      int(val),
			MoodTag:    key,
			Percentage: (val / float64(len(tags))) * 100.0,
		}
		moodTagFrequencies = append(moodTagFrequencies, mtf)
	}

	slices.SortFunc(moodTagFrequencies, func(a, b models.MoodTagFrequency) int {
		if a.Percentage > b.Percentage {
			return -1
		} else if a.Percentage < b.Percentage {
			return 1
		} else {
			return 0
		}
	})

	return moodTagFrequencies
}
