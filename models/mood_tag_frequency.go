package models

type MoodTagFrequency struct {
	MoodTag    string  `json:"moodTag"`
	Count      int     `json:"count"`
	Percentage float32 `json:"percentage"`
}
