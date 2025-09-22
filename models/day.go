package models

type Day struct {
	Date           string         `json:"date"`
	AvgRating      float64            `json:"avgRating"`
	JournalEntries []JournalEntry `json:"JournalEntries"`
}
