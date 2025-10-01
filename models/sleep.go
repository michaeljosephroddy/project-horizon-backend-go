package models

type Sleep struct {
	UserID        string  `json:"userId"`
	Granularity   string  `json:"granularity"`
	StartDate     string  `json:"startDate"`
	EndDate       string  `json:"endDate"`
	MovingAvg     float64 `json:"movingAvg"`
	Trend         string  `json:"sleepTrend"`
	StdDeviation  float64 `json:"stdDeviation"`
	AvgHoursSlept float64 `json:"avgHoursSlept"`
	Stability     float64 `json:"stability"`
}
