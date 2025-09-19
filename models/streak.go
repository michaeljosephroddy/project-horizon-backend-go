package models

type Streak struct {
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	StreakType string `json:"streakType"`
	Days []Day `json:"days"`
	NumDays int `json:"numDays"`
}
