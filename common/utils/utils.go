// common utils
package common_utils

import (
	"fmt"
	"regexp"
	"time"
)

func MatchURL(pattern string, url string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(url)
}

func GetUserIDFromPath(path string) (string, error) {
	re := regexp.MustCompile(`/users/([^/]+)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 || matches[1] == "" {
		return "", fmt.Errorf("no user ID found in path")
	}

	return matches[1], nil
}
func ParseDates(a string, b string) (time.Time, time.Time) {
	const layout = "2006-01-02"
	aParsed, _ := time.Parse(layout, a)
	bParsed, _ := time.Parse(layout, b)
	return aParsed, bParsed
}
