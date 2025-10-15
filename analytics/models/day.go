package models

import "time"

type Day struct {
	Date           time.Time       `json:"date"`
	DailyAvgRating float64         `json:"dailyAvgRating"`
	MoodLogs       []MoodLog       `json:"moodLogs"`
	SleepLogs      []SleepLog      `json:"sleepLogs"`
	MedicationLogs []MedicationLog `json:"medicationLogs"`
}
