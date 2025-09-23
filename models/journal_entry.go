package models

type JournalEntry struct {
	JournalEntryID int      `json:"journalEntryId"`
	UserID         string   `json:"userId"`
	MoodRating     int      `json:"moodRating"`
	Note           string   `json:"note"`
	CreatedAt      string   `json:"createdAt"`
	MoodTags       []string `json:"moodTags"`
}
