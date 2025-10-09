package models

type SleepLog struct {
	SleepLogID      string  `json:"sleepLogId"`
	UserID          string  `json:"userId"`
	HoursSlept      float64 `json:"hoursSlept"`
	SleepQualityTag string  `json:"sleepQualityTag"`
	Note            string  `json:"note"`
	SleepDate       string  `json:"sleepDate"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}
