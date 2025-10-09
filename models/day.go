package models

type Day struct {
	Date           string     `json:"date"`
	DailyAvgRating float64    `json:"dailyAvgRating"`
	MoodLogs       []MoodLog  `json:"moodLogs"`
	SleepLogs      []SleepLog `json:"sleepLogs"`
}
