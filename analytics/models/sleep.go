package models

import "time"

// Sleep Log
type SleepLog struct {
	SleepLogID      string    `json:"sleepLogId"`
	UserID          string    `json:"userId"`
	HoursSlept      float64   `json:"hoursSlept"`
	SleepQualityTag string    `json:"sleepQualityTag"`
	Note            string    `json:"note"`
	SleepDate       time.Time `json:"sleepDate"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// Sleep Metric
type SleepMetric struct {
	UserID               string    `json:"userID"`
	Granularity          string    `json:"granularity"`
	StartDate            time.Time `json:"startDate"`
	EndDate              time.Time `json:"endDate"`
	AvgSleepHours        float64   `json:"avgSleepHours"`
	MovingAvg            float64   `json:"movingAvg"`
	SleepTrend           string    `json:"sleepTrend"`
	StdDeviation         float64   `json:"stdDeviation"`
	Stability            string    `json:"stability"`
	SleepQualityTagStats []TagStat `json:"sleepQualityTagStats"`
	SleepDiffs           SleepDiff `json:"sleepDiffs"`
}

// Sleep Diff
type SleepDiff struct {
	AvgSleepHours MetricChange `json:"avgSleepHours"`
	Trend         ShiftChange       `json:"trend"`
	Stability     MetricChange `json:"stability"`
	TopQualityTag MetricChange `json:"topQualityTag"`
}

type ShiftChange struct {
	Description string  `json:"description"` // "increasing → stable"
	Change      float64 `json:"change"`
}
