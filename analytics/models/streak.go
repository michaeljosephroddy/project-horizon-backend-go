package models

import "time"

type Streak struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	NumDays   int       `json:"numDays"`
	Days      []Day     `json:"days"`
}
