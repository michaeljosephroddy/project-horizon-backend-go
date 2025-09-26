package models

type Metrics struct {
	UserID               string             `json:"userId"`
	Granularity          string             `json:"granularity"`
	PeriodStart          string             `json:"periodStart"`
	PeriodEnd            string             `json:"periodEnd"`
	Trend                string             `json:"moodTrend"`
	Stability            string             `json:"moodStability"`
	AvgMoodRatingPeriod  float64            `json:"avgMoodRatingPeriod"`
	TopMoodsPeriod       []MoodTagFrequency `json:"topMoodsPeriod"`
	TopMoodsPositiveDays []MoodTagFrequency `json:"topMoodsPositiveDays"`
	TopMoodsNegativeDays []MoodTagFrequency `json:"topMoodsNegativeDays"`
	PositiveStreaks      []Streak           `json:"positiveStreaks"`
	NegativeStreaks      []Streak           `json:"negativeStreaks"`
	PositiveDays         []Day              `json:"positiveDays"`
	NegativeDays         []Day              `json:"negativeDays"`
	PeriodDiffs          PeriodDiff         `json:"periodDiffs"`
}
