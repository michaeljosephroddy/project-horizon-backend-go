package database

import (
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
)

type SleepLogRepository struct {
	db *sql.DB
}

func NewSleepLogRepository(dbConnection *sql.DB) *SleepLogRepository {
	return &SleepLogRepository{
		db: dbConnection,
	}
}

func (slr *SleepLogRepository) AvgSleepHours(userID string, startDate time.Time, endDate time.Time) float64 {

	rows, queryErr := slr.db.Query(avgSleepHoursQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var avgSleepHours sql.NullFloat64

	for rows.Next() {
		scanErr := rows.Scan(&avgSleepHours)
		if scanErr != nil {
			panic(scanErr)
		}
	}

	if !avgSleepHours.Valid {
		return 0
	}

	return avgSleepHours.Float64
}

func (slr *SleepLogRepository) MovingAvgSleep(userID string, startDate time.Time, endDate time.Time, numDaysPreceding string) []models.MovingAverage {

	query := fmt.Sprintf(sleepMovingAvgQuery, numDaysPreceding)
	rows, queryErr := slr.db.Query(query, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var movingAvg models.MovingAverage
	var movingAverages []models.MovingAverage

	for rows.Next() {
		var date sql.NullString
		var movingAvgVal sql.NullFloat64

		scanErr := rows.Scan(
			&date,
			&movingAvgVal,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		if valid := date.Valid; valid {
			layout := "2006-01-02"
			d, _ := time.Parse(layout, date.String)
			movingAvg.Date = d
		}

		if valid := movingAvgVal.Valid; valid {
			movingAvg.MovingAvg = movingAvgVal.Float64
		}

		movingAverages = append(movingAverages, movingAvg)
	}

	if movingAverages == nil {
		return make([]models.MovingAverage, 0)
	}

	return movingAverages
}

func (slr *SleepLogRepository) StandardDeviation(userID string, startDate time.Time, endDate time.Time) float64 {

	rows, queryErr := slr.db.Query(sleepStdDevQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var standardDeviation sql.NullFloat64

	for rows.Next() {
		scanErr := rows.Scan(&standardDeviation)
		if scanErr != nil {
			panic(scanErr)
		}
	}

	if !standardDeviation.Valid {
		return 0.0
	}

	return standardDeviation.Float64
}

func (slr *SleepLogRepository) SleepQualityTagStat(userID string, startDate time.Time, endDate time.Time) []models.TagStat {

	rows, queryErr := slr.db.Query(sleepQualityTagStatQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var sleepTagStat models.TagStat
	var sleepTagFrequencies []models.TagStat

	for rows.Next() {
		scanErr := rows.Scan(
			&sleepTagStat.TagName,
			&sleepTagStat.Count,
			&sleepTagStat.Percentage,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		sleepTagFrequencies = append(sleepTagFrequencies, sleepTagStat)
	}

	if sleepTagFrequencies == nil {
		return make([]models.TagStat, 0)
	}

	slices.SortFunc(sleepTagFrequencies, func(a, b models.TagStat) int {
		return int(b.Percentage - a.Percentage)
	})

	return sleepTagFrequencies

}

func (slr *SleepLogRepository) SleepLogs(userID string, startDate time.Time, endDate time.Time) []models.SleepLog {
	rows, queryErr := slr.db.Query(sleepLogsQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var sleepLogs []models.SleepLog
	for rows.Next() {
		var sleepLog models.SleepLog
		var sleepDateStr, createdAtStr, updatedAtStr string

		scanErr := rows.Scan(
			&sleepLog.SleepLogID,
			&sleepLog.UserID,
			&sleepLog.HoursSlept,
			&sleepLog.SleepQualityTag,
			&sleepLog.Note,
			&sleepDateStr,
			&createdAtStr,
			&updatedAtStr,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		// Parse the date strings into time.Time
		parsedSleepDate, err := time.Parse("2006-01-02", sleepDateStr)
		if err != nil {
			panic(err)
		}
		parsedCreatedAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			panic(err)
		}
		parsedUpdatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			panic(err)
		}

		sleepLog.SleepDate = parsedSleepDate
		sleepLog.CreatedAt = parsedCreatedAt
		sleepLog.UpdatedAt = parsedUpdatedAt

		sleepLogs = append(sleepLogs, sleepLog)
	}

	if sleepLogs == nil {
		return make([]models.SleepLog, 0)
	}
	return sleepLogs
}
