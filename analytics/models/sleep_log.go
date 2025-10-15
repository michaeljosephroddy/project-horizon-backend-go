package models

import "time"

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
