// common utils
package utils

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

func ParseDate(a string) time.Time {
	const layout = "2006-01-02"
	aParsed, _ := time.Parse(layout, a)
	return aParsed
}
