package models

import "time"

// Medication Log
type MedicationLog struct {
	MedicationLogID int          `json:"medicationLogId"`
	UserID          int          `json:"userId"`
	TakenAt         time.Time    `json:"takenAt"` // date stored as string
	Note            string       `json:"note"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	Medications     []Medication `json:"medications"`
}

// Medication
type Medication struct {
	MedicationID int    `json:"medicationId"`
	Name         string `json:"name"`
	Dosage       string `json:"dosage"`
}

// TODO finish this
type MedicationMetric struct {
	UserID         string          `json:"userID"`
	Granularity    string          `json:"granularity"`
	StartDate      time.Time       `json:"startDate"`
	EndDate        time.Time       `json:"endDate"`
	CurrentMeds    []Medication    `json:"currentMeds"`
	MedicationLogs []MedicationLog `json:"medicationLogs"`
}
