package models

type Medication struct {
	MedicationID int    `json:"medicationId"`
	Name         string `json:"name"`
	Dosage       string `json:"dosage"`
}
