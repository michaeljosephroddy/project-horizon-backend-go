package database

import (
	"database/sql"
	"encoding/json"
	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type MedicationLogRepository struct {
	db *sql.DB
}

func NewMedicationLogRepository(dbConnection *sql.DB) *MedicationLogRepository {
	return &MedicationLogRepository{
		db: dbConnection,
	}
}

func (mlr *MedicationLogRepository) MedicationLogs(userID string, startDate string, endDate string) []models.MedicationLog {
	rows, queryErr := mlr.db.Query(medicationLogQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var medicationLogs []models.MedicationLog

	for rows.Next() {
		var medicationLog models.MedicationLog
		var medicationsJSON string // holds the JSON array from SQL

		scanErr := rows.Scan(
			&medicationLog.MedicationLogID,
			&medicationLog.UserID,
			&medicationLog.TakenAt,
			&medicationLog.Notes,
			&medicationLog.CreatedAt,
			&medicationLog.UpdatedAt,
			&medicationsJSON,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		// Parse JSON into Go struct
		var meds []models.Medication
		if err := json.Unmarshal([]byte(medicationsJSON), &meds); err != nil {
			panic(err)
		}
		medicationLog.Medications = meds

		medicationLogs = append(medicationLogs, medicationLog)
	}

	return medicationLogs
}
