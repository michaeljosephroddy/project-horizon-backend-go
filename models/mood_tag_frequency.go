package models

type MoodTagFrequency struct {
	MoodTagID string `json:"moodTagId"`
	Name string `json:"name"`
	Count int `json:"count"`
	Percentage float32 `json:"percentage"`
}
