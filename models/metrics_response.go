package models

type MetricsResponse struct {
	Trend string `json:"moodTrend"`
	Stability string `json:"moodStability"`
	MoodTagFrequency []MoodTagFrequency `json:"moodTagFrequencies"`
	PositiveStreaks []Streak `json:"positiveStreaks"`
	NegativeStreaks []Streak `json:"negativeStreaks"`
	PositiveDays []Day `json:"positiveDays"`
	NegativeDays []Day `json:"negativeDays"`
	// TODO add in the other data points best day, worst day etc..
}
