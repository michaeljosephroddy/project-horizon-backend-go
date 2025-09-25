package models

type MetricsResponse struct {
	CurrentPeriod      Period     `json:"currentPeriod"`
	PreviousPeriodDiff PeriodDiff `json:"previousPeriodDiff"`
}
