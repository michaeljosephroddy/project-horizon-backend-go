package models

type Streak struct {
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	NumDays int `json:"numDays"`
}
