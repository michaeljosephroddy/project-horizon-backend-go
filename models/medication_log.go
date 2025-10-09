package models

type MedicationLog struct {
	MedicationLogID int    `json:"medicationLogId"`
	UserID          int    `json:"userId"`
	MedicationID    int    `json:"medicationId"`
	TakenAt         string `json:"takenAt"` // date stored as string
	Taken           bool   `json:"taken"`
	Dosage          string `json:"dosage"`
	Notes           string `json:"notes"`
}
