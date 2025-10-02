package database

import (
	"database/sql"
	"fmt"

	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type SleepLogRepository struct {
	db *sql.DB
}

func NewSleepLogRepository(dbConnection *sql.DB) *SleepLogRepository {
	return &SleepLogRepository{
		db: dbConnection,
	}
}

func (slr *SleepLogRepository) AvgSleepHours(userID string, startDate string, endDate string) float64 {

	rows, queryErr := slr.db.Query(avgSleepHoursQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}

	var avgSleepHours float64
	if next := rows.Next(); next {
		rows.Scan(&avgSleepHours)
	}

	return avgSleepHours
}

func (slr *SleepLogRepository) MovingAvgSleep(userID string, startDate string, endDate string, numDaysPreceding string) []models.MovingAverage {

	query := fmt.Sprintf(sleepMovingAvgQuery, numDaysPreceding)
	rows, queryErr := slr.db.Query(query, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}

	var movingAvg models.MovingAverage
	var movingAverages []models.MovingAverage

	for rows.Next() {

		scanErr := rows.Scan(
			&movingAvg.Date,
			&movingAvg.MovingAvg,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		movingAverages = append(movingAverages, movingAvg)
	}

	if movingAverages == nil {
		return make([]models.MovingAverage, 0)
	}

	return movingAverages
}

func (slr *SleepLogRepository) StandardDeviation(userID string, startDate string, endDate string) float64 {

	rows, queryErr := slr.db.Query(sleepStdDevQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var standardDeviation sql.NullFloat64

	if next := rows.Next(); next {
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

/* func (slr *SleepLogRepository) SleepQualityTagFrequency(userID string, startDate string, endDate string) []models.TagFrequency {

	rows, queryErr := slr.db.Query(sleepQualityTagFrequencyQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}

} */
