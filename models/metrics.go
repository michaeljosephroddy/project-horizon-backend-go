package models

type Metrics struct {
	UserID               string             `json:"userId"`
	Granularity          string             `json:"granularity"`
	StartDate            string             `json:"startDate"`
	EndDate              string             `json:"endDate"`
	MovingAvg            float64            `json:"movingAvg"`
	Trend                string             `json:"moodTrend"`
	StdDeviation         float64            `json:"stdDeviation"`
	Stability            string             `json:"moodStability"`
	AvgMoodRating        float64            `json:"avgMoodRating"`
	TopMoods             []MoodTagFrequency `json:"topMoods"`
	TopMoodsPositiveDays []MoodTagFrequency `json:"topMoodsPositiveDays"`
	TopMoodsNeutralDays  []MoodTagFrequency `json:"topMoodsNeutralDays"`
	TopMoodsNegativeDays []MoodTagFrequency `json:"topMoodsNegativeDays"`
	TopMoodsClinicalDays []MoodTagFrequency `json:"topMoodsClinicalDays"`
	PositiveStreaks      []Streak           `json:"positiveStreaks"`
	NeutralStreaks       []Streak           `json:"neutralStreaks"`
	NegativeStreaks      []Streak           `json:"negativeStreaks"`
	ClinicalStreaks      []Streak           `json:"clinicalStreaks"`
	PositiveDays         []Day              `json:"positiveDays"`
	NeutralDays          []Day              `json:"neutralDays"`
	NegativeDays         []Day              `json:"negativeDays"`
	ClinicalDays         []Day              `json:"clinicalDays"`
	Diffs                Diff               `json:"diffs"`
}
