package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
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

// OverviewStats gets high-level medication statistics
func (mlr *MedicationLogRepository) OverviewStats(userID string, startDate, endDate time.Time) (int, float64, float64, float64, error) {
	query := `
		SELECT 
			COUNT(*) as total_logs,
			COUNT(DISTINCT DATE(ml.taken_at)) as days_with_logs,
			DATEDIFF(?, ?) + 1 as total_days,
			AVG(med_count) as avg_meds_per_log
		FROM medication_log ml
		LEFT JOIN (
			SELECT medication_log_id, COUNT(*) as med_count
			FROM medication_log_item
			GROUP BY medication_log_id
		) med_counts ON ml.medication_log_id = med_counts.medication_log_id
		LEFT JOIN medication_log_item mli ON ml.medication_log_id = mli.medication_log_id
		WHERE ml.user_id = ?
			AND DATE(ml.taken_at) BETWEEN ? AND ?
	`

	var totalLogs, daysWithLogs, totalDays int
	var avgMedsPerLog sql.NullFloat64

	err := mlr.db.QueryRow(query, endDate, startDate, userID, startDate, endDate).Scan(
		&totalLogs, &daysWithLogs, &totalDays, &avgMedsPerLog,
	)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get overview stats: %w", err)
	}

	adherenceRate := 0.0
	if totalDays > 0 {
		adherenceRate = (float64(daysWithLogs) / float64(totalDays)) * 100
	}

	avgLogsPerDay := 0.0
	if totalDays > 0 {
		avgLogsPerDay = float64(totalLogs) / float64(totalDays)
	}

	avgMedsPerLogValue := 0.0
	if avgMedsPerLog.Valid {
		avgMedsPerLogValue = avgMedsPerLog.Float64
	}

	return totalLogs, adherenceRate, avgLogsPerDay, avgMedsPerLogValue, nil
}

// MedicationDetailedStats returns comprehensive stats for each medication
func (mlr *MedicationLogRepository) MedicationDetailedStats(userID string, startDate, endDate time.Time) ([]models.MedicationStats, error) {
	// First, get basic stats per medication
	query := `SELECT 
			m.medication_id,
			m.name,
			COUNT(*) as total_doses,
			COUNT(DISTINCT DATE(ml.taken_at)) as days_active,
			DATEDIFF(?, ?) + 1 as total_days
		FROM medication_log ml
		JOIN medication_log_item mli ON ml.medication_log_id = mli.medication_log_id
		JOIN medication m ON mli.medication_id = m.medication_id
		WHERE ml.user_id = ?
			AND DATE(ml.taken_at) BETWEEN ? AND ?
		GROUP BY m.medication_id, m.name
		ORDER BY total_doses DESC`

	rows, err := mlr.db.Query(query, endDate, startDate, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication stats: %w", err)
	}
	defer rows.Close()

	var stats []models.MedicationStats

	for rows.Next() {
		var medID, totalDoses, daysActive, totalDays int
		var name string

		if err := rows.Scan(&medID, &name, &totalDoses, &daysActive, &totalDays); err != nil {
			return nil, fmt.Errorf("failed to scan medication stat: %w", err)
		}

		stat := models.MedicationStats{
			MedicationID:   medID,
			Name:           name,
			TotalDoses:     totalDoses,
			DaysActive:     daysActive,
			AvgDosesPerDay: float64(totalDoses) / float64(totalDays),
		}

		// Get timing stats
		timingStats, err := mlr.getTimingStats(userID, medID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		stat.AvgTakenAtTime = timingStats.AvgTime
		stat.TimingStdDevMinutes = timingStats.StdDevMinutes
		stat.TimingDescription = timingStats.Description
		stat.EarliestTime = timingStats.EarliestTime
		stat.LatestTime = timingStats.LatestTime

		// Get streaks
		longestStreak, currentStreak, err := mlr.getStreaks(userID, medID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		stat.LongestStreak = longestStreak
		stat.CurrentStreak = currentStreak

		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

func (mlr *MedicationLogRepository) getTimingStats(userID string, medID int, startDate, endDate time.Time) (models.TimingStats, error) {
	query := `SELECT 
				TIME_FORMAT(SEC_TO_TIME(AVG(TIME_TO_SEC(TIME(ml.taken_at)))), '%H:%i:%s') as avg_time,
				STD(TIME_TO_SEC(TIME(ml.taken_at))) / 60 as std_dev_minutes,
				TIME_FORMAT(MIN(TIME(ml.taken_at)), '%H:%i:%s') as earliest_time,
				TIME_FORMAT(MAX(TIME(ml.taken_at)), '%H:%i:%s') as latest_time
			FROM medication_log ml
			JOIN medication_log_item mli ON ml.medication_log_id = mli.medication_log_id
			WHERE ml.user_id = ?
				AND mli.medication_id = ?
				AND DATE(ml.taken_at) BETWEEN ? AND ?`

	var avgTime, earliestTime, latestTime string
	var stdDevMinutes sql.NullFloat64

	err := mlr.db.QueryRow(query, userID, medID, startDate, endDate).Scan(
		&avgTime, &stdDevMinutes, &earliestTime, &latestTime,
	)
	if err != nil {
		return models.TimingStats{}, fmt.Errorf("failed to get timing stats: %w", err)
	}

	stdDev := 0.0
	if stdDevMinutes.Valid {
		stdDev = stdDevMinutes.Float64
	}

	// Format description like "8:47 AM ± 45 minutes"
	description := formatTimingDescription(avgTime, stdDev)

	return models.TimingStats{
		AvgTime:       avgTime,
		StdDevMinutes: stdDev,
		Description:   description,
		EarliestTime:  earliestTime,
		LatestTime:    latestTime,
	}, nil
}

func formatTimingDescription(avgTime string, stdDevMinutes float64) string {
	// Parse avgTime (HH:MM:SS format)
	t, err := time.Parse("15:04:05", avgTime)
	if err != nil {
		return avgTime
	}

	// Format to 12-hour with AM/PM
	timeStr := t.Format("3:04 PM")

	// Round std dev to nearest minute
	stdDevRounded := int(math.Round(stdDevMinutes))

	return fmt.Sprintf("%s ± %d minutes", timeStr, stdDevRounded)
}

func (mlr *MedicationLogRepository) getStreaks(userID string, medID int, startDate, endDate time.Time) (int, int, error) {
	query := `
		SELECT DISTINCT DATE(ml.taken_at) as log_date
		FROM medication_log ml
		JOIN medication_log_item mli ON ml.medication_log_id = mli.medication_log_id
		WHERE ml.user_id = ?
			AND mli.medication_id = ?
			AND DATE(ml.taken_at) BETWEEN ? AND ?
		ORDER BY log_date
	`

	rows, err := mlr.db.Query(query, userID, medID, startDate, endDate)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get streaks: %w", err)
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var dateStr string // Changed to string
		if err := rows.Scan(&dateStr); err != nil {
			return 0, 0, fmt.Errorf("failed to scan date: %w", err)
		}

		// Parse date string
		const dateLayout = "2006-01-02"
		date, err := time.Parse(dateLayout, dateStr)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse date: %w", err)
		}

		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0, 0, nil
	}

	longestStreak := 1
	currentStreak := 1
	tempStreak := 1

	for i := 1; i < len(dates); i++ {
		daysDiff := int(dates[i].Sub(dates[i-1]).Hours() / 24)

		if daysDiff == 1 {
			tempStreak++
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
		} else {
			tempStreak = 1
		}
	}

	// Calculate current streak (from end date backwards)
	if len(dates) > 0 {
		lastDate := dates[len(dates)-1]
		daysSinceLastLog := int(endDate.Sub(lastDate).Hours() / 24)

		if daysSinceLastLog <= 1 {
			currentStreak = 1
			for i := len(dates) - 2; i >= 0; i-- {
				daysDiff := int(dates[i+1].Sub(dates[i]).Hours() / 24)
				if daysDiff == 1 {
					currentStreak++
				} else {
					break
				}
			}
		} else {
			currentStreak = 0
		}
	}

	return longestStreak, currentStreak, rows.Err()
}
