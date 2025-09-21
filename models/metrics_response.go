package models

type MetricsResponse struct {
	Trend string `json:"moodTrend"`
	Stability string `json:"moodStability"`
	MoodTagFrequency []MoodTagFrequency `json:"moodTagFrequencies"`
	HighStreaks []Streak `json:"highStreaks"`
	LowStreaks []Streak `json:"lowStreaks"`
	// TODO add in the other data points best day, worst day etc..
}
