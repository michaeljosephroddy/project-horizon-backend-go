package model

import "time"

type MovingAverage struct {
	Date      time.Time `json:"date"`
	MovingAvg float64   `json:"movingAvg"`
}

type Day struct {
	Date           time.Time       `json:"date"`
	DailyAvgRating float64         `json:"dailyAvgRating"`
	MoodLogs       []MoodLog       `json:"moodLogs"`
	SleepLogs      []SleepLog      `json:"sleepLogs"`
	MedicationLogs []MedicationLog `json:"medicationLogs"`
}

type Streak struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	NumDays   int       `json:"numDays"`
	Days      []Day     `json:"days"`
}

type TagStat struct {
	TagName    string  `json:"tagName"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type MetricChange struct {
	Current       float64 `json:"current"`
	Previous      float64 `json:"previous"`
	Change        float64 `json:"change"`
	PercentChange float64 `json:"percentChange"`
	Shift         string  `json:"shift"`
}

type TimingStats struct {
	AvgTime       string
	StdDevMinutes float64
	Description   string
	EarliestTime  string
	LatestTime    string
}
