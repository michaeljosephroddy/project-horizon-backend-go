package models

type MetricsResponse struct {
	Trend string `json:"moodTrend"`
	MoodTagFrequency []MoodTagFrequency `json:"moodTagFrequencies"`
	Stability int `json:"moodStability"`
	HighStreak []Day `json:"highStreak"`
	LowStreak []Day `json:"lowStreak"`
	// TODO add in the other data points best day, worst day etc..
}
