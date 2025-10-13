package models

type MedicationLog struct {
	MedicationLogID int          `json:"medicationLogId"`
	UserID          int          `json:"userId"`
	TakenAt         string       `json:"takenAt"` // date stored as string
	Notes           string       `json:"notes"`
	CreatedAt       string       `json:"createdAt"`
	UpdatedAt       string       `json:"updatedAt"`
	Medications     []Medication `json:"medications"`
}
