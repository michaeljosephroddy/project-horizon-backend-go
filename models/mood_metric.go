package models

type MoodMetric struct {
	UserID                   string    `json:"userId"`
	Granularity              string    `json:"granularity"`
	StartDate                string    `json:"startDate"`
	EndDate                  string    `json:"endDate"`
	MovingAvg                float64   `json:"movingAvg"`
	MoodTrend                string    `json:"moodTrend"`
	StdDeviation             float64   `json:"stdDeviation"`
	Stability                string    `json:"moodStability"`
	AvgMoodRating            float64   `json:"avgMoodRating"`
	MoodTagStats             []TagStat `json:"moodTagStats"`
	MoodTagStatsPositiveDays []TagStat `json:"moodTagStatsPositiveDays"`
	MoodTagStatsNeutralDays  []TagStat `json:"moodTagStatsNeutralDays"`
	MoodTagStatsNegativeDays []TagStat `json:"moodTagStatsNegativeDays"`
	MoodTagStatsClinicalDays []TagStat `json:"moodTagStatsClinicalDays"`
	PositiveStreaks          []Streak  `json:"positiveStreaks"`
	NeutralStreaks           []Streak  `json:"neutralStreaks"`
	NegativeStreaks          []Streak  `json:"negativeStreaks"`
	ClinicalStreaks          []Streak  `json:"clinicalStreaks"`
	PositiveDays             []Day     `json:"positiveDays"`
	NeutralDays              []Day     `json:"neutralDays"`
	NegativeDays             []Day     `json:"negativeDays"`
	ClinicalDays             []Day     `json:"clinicalDays"`
	MoodDiffs                MoodDiff  `json:"moodDiffs"`
}
