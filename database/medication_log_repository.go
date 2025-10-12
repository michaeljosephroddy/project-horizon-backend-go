package database

import (
	"database/sql"

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

		scanErr := rows.Scan(
			&medicationLog.MedicationLogID,
			&medicationLog.UserID,
			&medicationLog.MedicationID,
			&medicationLog.TakenAt,
			&medicationLog.Dosage,
			&medicationLog.Notes,
		)

		if scanErr != nil {
			panic(scanErr)
		}

		medicationLogs = append(medicationLogs, medicationLog)
	}
	
	if medicationLogs != nil {
		return make([]models.MedicationLog, 0)
	}

	return medicationLogs
}
