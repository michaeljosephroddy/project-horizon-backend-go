package models

type SleepMetric struct {
	UserID              string         `json:"userId"`
	Granularity         string         `json:"granularity"`
	StartDate           string         `json:"startDate"`
	EndDate             string         `json:"endDate"`
	MovingAvg           float64        `json:"movingAvg"`
	SleepTrend          string         `json:"sleepTrend"`
	StdDeviation        float64        `json:"stdDeviation"`
	Stability           string         `json:"stability"`
	AvgSleepHours       float64        `json:"avgSleepHours"`
	TopSleepQualityTags []TagFrequency `json:"topSleepQualityTags"`
}
