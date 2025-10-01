package models

type Medication struct {
	UserID      string   `json:"userId"`
	Granularity string   `json:"granularity"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Medications []string `json:"medications"`
}
