package models

type Day struct {
	Date               string             `json:"date"`
	DailyAvgRating     float64            `json:"dailyAvgRating"`
	JournalEntries     []JournalEntry     `json:"journalEntries"`
	MoodTagFrequencies []MoodTagFrequency `json:"moodTagFrequencies"`
}
