package model

import "time"

type User struct {
	UserID       int       `json:"userId"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// UserMedicationDTO represents user medication with joined medication details
type UserMedicationDTO struct {
	UserMedicationID int       `json:"userMedicationId"`
	MedicationID     int       `json:"medicationId"`
	Name             string    `json:"name"`
	Dosage           *string   `json:"dosage"`
	StartDate        time.Time `json:"startDate"`
	Note             *string   `json:"note"`
	Description      *string   `json:"description"`
}
