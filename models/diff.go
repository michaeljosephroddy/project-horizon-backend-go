package models

type Diff struct {
	AvgMoodChange              float64 `json:"avgMoodChange"`   // +0.8
	TrendChange                string  `json:"trendChange"`     // "increasing → stable"
	StabilityChange            string  `json:"stabilityChange"` // "moderate → stable"
	VolatilityDelta            float64 `json:"volatilityDelta"` // -0.2
	VolatilityDeltaPercentage  float64 `json:volatilityDeltaPercentage`
	TopMoodShift               string  `json:"topMoodShift"`               // "HAPPY → SAD"
	TopMoodDelta               float64 `json:"topMoodDelta"`               // -12.0 (percentage points)
	TopPositiveMoodChange      string  `json:"topPositiveMoodChange"`      // "JOY +8%"
	TopNegativeMoodChange      string  `json:"topNegativeMoodChange"`      // "ANGER -5%"
	PositiveDaysDelta          int     `json:"positiveDaysDelta"`          // +3
	NegativeDaysDelta          int     `json:"negativeDaysDelta"`          // -2
	PositiveRatioChange        float64 `json:"positiveRatioChange"`        // +12.0
	LongestPositiveStreakDelta int     `json:"longestPositiveStreakDelta"` // +2
	LongestNegativeStreakDelta int     `json:"longestNegativeStreakDelta"` // -1
}
