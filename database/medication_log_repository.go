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

func (mlr *MedicationLogRepository) MedicationLogs(userID string, startDate, endDate time.Time) []models.MedicationLog {

	rows, queryErr := mlr.db.Query(medicationLogQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var medicationLogs []models.MedicationLog

	for rows.Next() {
		var medicationLog models.MedicationLog
		var medicationsJSON string
		var note sql.NullString
		var takenAtStr, createdAtStr, updatedAtStr string

		scanErr := rows.Scan(
			&medicationLog.MedicationLogID,
			&medicationLog.UserID,
			&takenAtStr,
			&note,
			&createdAtStr,
			&updatedAtStr,
			&medicationsJSON,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		// Parse taken_at
		takenAt, err := time.Parse("2006-01-02 15:04:05", takenAtStr)
		if err != nil {
			panic(fmt.Errorf("failed to parse taken_at: %v", err))
		}
		medicationLog.TakenAt = takenAt

		// Parse created_at
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			panic(fmt.Errorf("failed to parse created_at: %v", err))
		}
		medicationLog.CreatedAt = createdAt

		// Parse updated_at
		updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			panic(fmt.Errorf("failed to parse updated_at: %v", err))
		}
		medicationLog.UpdatedAt = updatedAt

		if note.Valid {
			medicationLog.Note = note.String
		}

		// Parse JSON field
		var meds []models.Medication
		if err := json.Unmarshal([]byte(medicationsJSON), &meds); err != nil {
			panic(err)
		}
		medicationLog.Medications = meds

		medicationLogs = append(medicationLogs, medicationLog)
	}

	return medicationLogs
}
