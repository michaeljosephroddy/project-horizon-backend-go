package utils

import (
	"regexp"
	"slices"
	"strings"
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

/* func MoodTagFrequencies(data []models.Day) []models.MoodTagFrequency {

	var freq map[string]interface{}

	for i := 0; i < len(data); i++ {
		for _, val := range data[i].MoodTagFrequencies {

			if val.Count > count {

			}
			freq["tag"] = val.MoodTag
			freq["count"] = val.Count
			freq["percentage"] = val.Percentage
		}
	}

} */
