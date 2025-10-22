package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
	"time"
)

type MedicationLogRepository struct {
	db *sql.DB
}

func NewMedicationLogRepository(dbConnection *sql.DB) *MedicationLogRepository {
	return &MedicationLogRepository{
		db: dbConnection,
	}
}

func (mlr *MedicationLogRepository) MedicationLogs(userID string, startDate, endDate time.Time) ([]models.MedicationLog, error) {
	rows, err := mlr.db.Query(medicationLogQuery, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query medication logs: %w", err)
	}
	defer rows.Close()

	var medicationLogs []models.MedicationLog

	for rows.Next() {
		var medicationLog models.MedicationLog
		var medicationsJSON string
		var note sql.NullString
		var takenAtStr, createdAtStr, updatedAtStr string

		err := rows.Scan(
			&medicationLog.MedicationLogID,
			&medicationLog.UserID,
			&takenAtStr,
			&note,
			&createdAtStr,
			&updatedAtStr,
			&medicationsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan medication log row: %w", err)
		}

		const dateTimeLayout = "2006-01-02 15:04:05"

		// Parse taken_at
		takenAt, err := time.Parse(dateTimeLayout, takenAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse taken_at: %w", err)
		}
		medicationLog.TakenAt = takenAt

		// Parse created_at
		createdAt, err := time.Parse(dateTimeLayout, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		medicationLog.CreatedAt = createdAt

		// Parse updated_at
		updatedAt, err := time.Parse(dateTimeLayout, updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}
		medicationLog.UpdatedAt = updatedAt

		if note.Valid {
			medicationLog.Note = note.String
		}

		// Parse JSON field
		var meds []models.Medication
		if err := json.Unmarshal([]byte(medicationsJSON), &meds); err != nil {
			return nil, fmt.Errorf("failed to unmarshal medications JSON: %w", err)
		}
		medicationLog.Medications = meds

		medicationLogs = append(medicationLogs, medicationLog)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating medication log rows: %w", err)
	}

	return medicationLogs, nil
}
