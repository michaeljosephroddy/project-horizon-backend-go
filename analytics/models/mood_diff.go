package models

type MoodDiff struct {
	// Overall metrics
	AvgRatingChange float64 `json:"avgRatingChange"`
	TrendShift      string  `json:"trendShift"` // "increasing → stable"
	MovingAvgChange float64 `json:"movingAvgChange"`
	StabilityShift  string  `json:"stabilityShift"` // "moderate → stable"
	StabilityChange float64 `json:"stabilityChange"`

	// Top mood changes
	TopTagShift  string  `json:"topTagShift"`  // "HAPPY → SAD", renamed for clarity
	TopTagChange float64 `json:"topTagChange"` // changed to float64

	// Category-specific changes (consolidated)
	Categories CategoryDiffs `json:"categories"`
}

type CategoryDiffs struct {
	Positive CategoryDiff `json:"positive"`
	Neutral  CategoryDiff `json:"neutral"`
	Negative CategoryDiff `json:"negative"`
	Clinical CategoryDiff `json:"clinical"`
}

type CategoryDiff struct {
	TopTagChange float64 `json:"topTagChange"`
	DaysChange   int     `json:"daysChange"`
	StreakChange int     `json:"streakChange"`
}
