package models

type Day struct {
	Date           string     `json:"date"`
	DailyAvgRating float64    `json:"dailyAvgRating"`
	MoodLogs       []MoodLog  `json:"moodLogs"`
	SleepLogs      []SleepLog `json:"sleepLogs"`
	MoodTagStats   []TagStat  `json:"moodTagStats"`
}
