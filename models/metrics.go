package models

type Metrics struct {
	UserID               string             `json:"userId"`
	Granularity          string             `json:"granularity"`
	StartDate            string             `json:"startDate"`
	EndDate              string             `json:"endDate"`
	Trend                string             `json:"moodTrend"`
	StdDeviation         float64            `json:"stdDeviation"`
	Stability            string             `json:"moodStability"`
	AvgMoodRating        float64            `json:"avgMoodRating"`
	TopMoods             []MoodTagFrequency `json:"topMoods"`
	TopMoodsPositiveDays []MoodTagFrequency `json:"topMoodsPositiveDays"`
	TopMoodsNegativeDays []MoodTagFrequency `json:"topMoodsNegativeDays"`
	PositiveStreaks      []Streak           `json:"positiveStreaks"`
	NegativeStreaks      []Streak           `json:"negativeStreaks"`
	PositiveDays         []Day              `json:"positiveDays"`
	NegativeDays         []Day              `json:"negativeDays"`
	Diffs                Diff               `json:"diffs"`
}
