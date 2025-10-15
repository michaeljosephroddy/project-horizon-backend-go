package models

import "time"

type SleepMetric struct {
	UserID               string    `json:"userId"`
	Granularity          string    `json:"granularity"`
	StartDate            time.Time `json:"startDate"`
	EndDate              time.Time `json:"endDate"`
	MovingAvg            float64   `json:"movingAvg"`
	SleepTrend           string    `json:"sleepTrend"`
	StdDeviation         float64   `json:"stdDeviation"`
	Stability            string    `json:"stability"`
	AvgSleepHours        float64   `json:"avgSleepHours"`
	SleepQualityTagStats []TagStat `json:"sleepQualityTagStats"`
	SleepDiffs           SleepDiff `json:"sleepDiffs"`
}
