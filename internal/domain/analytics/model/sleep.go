package model

import "time"

// Sleep Log
type SleepLog struct {
	SleepLogID      int       `json:"sleepLogId"`
	UserID          int       `json:"userId"`
	HoursSlept      float64   `json:"hoursSlept"`
	SleepQualityTag string    `json:"sleepQualityTag"`
	Note            string    `json:"note"`
	SleepDate       time.Time `json:"sleepDate"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// Sleep Metric
type SleepMetric struct {
	UserID             int       `json:"userID"`
	Granularity        string    `json:"granularity"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
	AvgSleepHours      float64   `json:"avgSleepHours"`
	MovingAvg          float64   `json:"movingAvg"`
	SleepTrend         string    `json:"sleepTrend"`
	StdDeviation       float64   `json:"stdDeviation"`
	Stability          string    `json:"stability"`
	BestSleepDay       string    `json:"bestSleepDay"`
	WorstSleepDay      string    `json:"worstSleepDay"`
	TopSleepQualityTag TagStat   `json:"topSleepQualityTag"`
}

type DayOfWeekSleepPattern struct {
	DayOfWeek     string  `json:"dayOfWeek"`     // "Monday", "Tuesday", etc.
	DayNumber     int     `json:"dayNumber"`     // 1=Sunday, 2=Monday, etc. (MySQL DAYOFWEEK)
	AvgSleepHours float64 `json:"avgSleepHours"` // Average hours slept on this day
	TotalEntries  int     `json:"totalEntries"`  // Number of logs for this day
}
