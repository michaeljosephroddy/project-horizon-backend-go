package models

type MoodTagFrequency struct {
	MoodTag    string  `json:"moodTag"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}
