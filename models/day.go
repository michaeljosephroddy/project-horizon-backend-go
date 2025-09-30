package models

type Day struct {
	Date               string             `json:"date"`
	DailyAvgRating     float64            `json:"dailyAvgRating"`
	MoodLogs           []MoodLog          `json:"moodLogs"`
	MoodTagFrequencies []MoodTagFrequency `json:"moodTagFrequencies"`
}
