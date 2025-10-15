package models

import "time"

type MedicationLog struct {
	MedicationLogID int          `json:"medicationLogId"`
	UserID          int          `json:"userId"`
	TakenAt         time.Time    `json:"takenAt"` // date stored as string
	Note            string       `json:"note"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	Medications     []Medication `json:"medications"`
}
