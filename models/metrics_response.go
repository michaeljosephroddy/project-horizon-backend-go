package models

type MetricsResponse struct {
	UserID             string             `json:"userId"`
	PeriodStart        string             `json:"periodStart"`
	PeriodEnd          string             `json:"periodEnd"`
	Trend              string             `json:"moodTrend"`
	Stability          string             `json:"moodStability"`
	MoodTagFrequencies []MoodTagFrequency `json:"moodTagFrquencies"`
	Top3Moods          []MoodTagFrequency `json:"top3MoodTags"`
	PositiveStreaks    []Streak           `json:"positiveStreaks"`
	NegativeStreaks    []Streak           `json:"negativeStreaks"`
	PositiveDays       []Day              `json:"positiveDays"`
	NegativeDays       []Day              `json:"negativeDays"`
	// TODO add in the other data points best day, worst day etc..
}
