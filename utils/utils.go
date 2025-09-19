package utils

import (
	"regexp"
	"slices"
	"strings"

	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
	"strconv"
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
