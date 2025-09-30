package models

type Diff struct {
	AvgMoodPercentChange             float64 `json:"avgMoodPercentChange"` // +0.8
	TrendShift                       string  `json:"trendShift"`           // "increasing → stable"
	MovingAvgPercentChange           float64 `json:"movingAvgPercentChange"`
	StabilityShift                   string  `json:"stabilityShift"` // "moderate → stable"
	StabilityPercentChange           float64 `json:"stabilityPercentChange"`
	TopMoodShift                     string  `json:"topMoodShift"`                     // "HAPPY → SAD"
	TopMoodPercentChange             string  `json:"topMoodPercentChange"`             // -12.0 (percentage points)
	TopMoodPositiveDaysPercentChange string  `json:"topMoodPositiveDaysPercentChange"` // "JOY +8%"
	TopMoodNeutralDaysPercentChange  string  `json:"topMoodNeutralDaysPercentChange"`
	TopMoodNegativeDaysPercentChange string  `json:"topMoodNegativeDaysPercentChange"` // "ANGER -5%"
	TopMoodClinicalDaysPercentChange string  `josn:"topMoodClinicalDaysPercentChange"`
	PositiveDaysChange               int     `json:"positiveDaysChange"` // +3
	NeutralDaysChange                int     `json:"neutralDaysChange"`
	NegativeDaysChange               int     `json:"negativeDaysChange"` // -2
	ClinicalDaysChange               int     `json:"clinicalDaysChange"`
	LongestPositiveStreakChange      int     `json:"longestPositiveStreakChange"` // +2
	LongestNeutralStreakChange       int     `json:"longestNeutralStreakChange"`
	LongestNegativeStreakChange      int     `json:"longestNegativeStreakChange"` // -1
	LongestClinicalStreakChange      int     `json:"longestClinicalStreakChange"`
}
