package models

import "time"

// Mood Metric
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

// Mood Diff
type MoodDiff struct {
	AvgRating  MetricChange `json:"avgRating"`
	Trend      ShiftChange      `json:"trend"`
	Stability  MetricChange `json:"stability"`
	TopTag     MetricChange `json:"topTag"`
	Categories CategoryDiffs    `json:"categories"`
}


type CategoryDiffs struct {
	Positive CategoryDiff `json:"positive"`
	Neutral  CategoryDiff `json:"neutral"`
	Negative CategoryDiff `json:"negative"`
	Clinical CategoryDiff `json:"clinical"`
}

type CategoryDiff struct {
	TopTag       MetricChange `json:"topTag"`
	DaysChange   int          `json:"daysChange"`
	StreakChange int          `json:"streakChange"`
}

// Mood Log
type MoodLog struct {
	MoodLogID  int       `json:"moodLogId"`
	UserID     string    `json:"userId"`
	MoodRating int       `json:"moodRating"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"createdAt"`
	MoodTags   []string  `json:"moodTags"`
}
