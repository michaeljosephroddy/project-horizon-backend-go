package models

type PeriodDiff struct {
	AvgMoodChange              float64 // +0.8
	TrendChange                string  // "increasing → stable"
	StabilityChange            string  // "moderate → stable"
	VolatilityDelta            float64 // -0.2
	TopMoodShift               string  // "HAPPY → SAD"
	TopMoodDelta               float64 // -12.0 (percentage points)
	TopPositiveMoodChange      string  // "JOY +8%"
	TopNegativeMoodChange      string  // "ANGER -5%"
	PositiveDaysDelta          int     // +3
	NegativeDaysDelta          int     // -2
	PositiveRatioChange        float64 // +12.0
	LongestPositiveStreakDelta int     // +2
	LongestNegativeStreakDelta int     // -1
}
