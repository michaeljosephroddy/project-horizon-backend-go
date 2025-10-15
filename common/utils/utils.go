// common utils
package common_utils

import (
	"regexp"
	"slices"
	"strings"
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

func ParseDates(a string, b string) (time.Time, time.Time) {
	const layout = "2006-01-02"
	aParsed, _ := time.Parse(layout, a)
	bParsed, _ := time.Parse(layout, b)
	return aParsed, bParsed
}
