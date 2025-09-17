package utils

import (
	"fmt"
	"regexp"
	"strings"
	"slices"
)

func MatchURL(pattern string, url string) bool {
	fmt.Println(pattern)
	fmt.Println(url)
	re := regexp.MustCompile(pattern)
	if re.MatchString(url) {
		fmt.Println("Matched!")
		return true
	} else {
		fmt.Println("Not matched!")
		return false
	}
}

func GetUserIdFromPath(path string) string {
	splitPath := strings.Split(path, "/")
	userIdIndex := slices.Index(splitPath, "users") + 1
	fmt.Println(splitPath[userIdIndex])
	return splitPath[userIdIndex]
}
