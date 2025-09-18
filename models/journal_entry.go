package models

type JournalEntry struct {
	JournalEntryID int `json:"journalEntryId"`
	UserID int `json:"userId"`
	MoodRating int `json:"moodRating"`
	Note string `json:"note"`
	CreatedAt string `json:"createdAt"`
}
