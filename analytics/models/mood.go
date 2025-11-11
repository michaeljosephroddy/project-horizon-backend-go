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
	TagStats             []TagStat `json:"tagStats"`
	PositiveDaysTagStats []TagStat `json:"positiveDaysTagStats"`
	NegativeDaysTagStats []TagStat `json:"negativeDaysTagStats"`
	NeutralDaysTagStats  []TagStat `json:"neutralDaysTagStats"`
	ClinicalDaysTagStats []TagStat `json:"clinicalDaysTagStats"`
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
