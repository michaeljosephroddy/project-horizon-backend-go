package models

type SleepDiff struct {
	AvgSleepHoursPercentChange      float64 `json:"avgSleepHoursPercentChange"`
	TrendShift                      string  `json:"trendShift"` // "increasing → stable"
	MovingAvgPercentChange          float64 `json:"movingAvgPercentChange"`
	StabilityShift                  string  `json:"stabilityShift"` // "moderate → stable"
	StabilityPercentChange          float64 `json:"stabilityPercentChange"`
	TopSleepQualityTagShift         string  `json:"topSleepQualityTagShift"`
	TopSleepQualityTagPercentChange string  `json:"topSleepQualityTagPercentChange"`
}
