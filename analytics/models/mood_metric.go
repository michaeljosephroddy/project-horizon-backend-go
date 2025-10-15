package models

import "time"

type MoodMetric struct {
	UserID      string    `json:"userId"`
	Granularity string    `json:"granularity"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`

	// Statistical measures
	MovingAvg    float64 `json:"movingAvg"`
	AvgRating    float64 `json:"avgRating"`
	Trend        string  `json:"trend"`
	StdDeviation float64 `json:"stdDeviation"`
	Stability    string  `json:"stability"`

	// Overall tag statistics
	TagStats []TagStat `json:"tagStats"`

	// Category-specific data (consolidated)
	Categories MoodCategories `json:"categories"`

	// Diffs
	Diffs MoodDiff `json:"diffs"`
}

type MoodCategories struct {
	Positive CategoryData `json:"positive"`
	Neutral  CategoryData `json:"neutral"`
	Negative CategoryData `json:"negative"`
	Clinical CategoryData `json:"clinical"`
}

type CategoryData struct {
	TagStats []TagStat `json:"tagStats"`
	Streaks  []Streak  `json:"streaks"`
	Days     []Day     `json:"days"`
}
