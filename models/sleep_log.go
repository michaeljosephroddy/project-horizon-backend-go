package models

type SleepLog struct {
	SleepLogID        string `json:"sleepLogId"`
	UserID            string `json:"userId"`
	HoursSlept        int    `json:"hoursSlept"`
	SleepQualityTagID string `json:"sleepQualityTagId"`
	Notes             string `json:"notes"`
	SleepDate         string `json:"sleepDate"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}
