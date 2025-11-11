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

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
)

func (slr *SleepLogRepository) AvgSleepHours(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	query := `SELECT Avg(hours_slept) AS avg_sleep_hours
FROM   sleep_log
WHERE  user_id = ?
       AND sleep_date BETWEEN ? AND ?;`

	rows, err := slr.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("failed to query average sleep hours: %w", err)
	}
	defer rows.Close()

	var avgSleepHours sql.NullFloat64

	for rows.Next() {
		err := rows.Scan(&avgSleepHours)
		if err != nil {
			return 0, fmt.Errorf("failed to scan average sleep hours: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating average sleep hours rows: %w", err)
	}

	if !avgSleepHours.Valid {
		return 0, nil
	}

	return avgSleepHours.Float64, nil
}

func (slr *SleepLogRepository) MovingAvgSleep(userID string, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]models.MovingAverage, error) {
	query := fmt.Sprintf(`WITH first_query
     AS (SELECT sleep_date AS DATE,
                Avg(hours_slept) AS avg_sleep_hours 
         FROM   sleep_log
         WHERE  user_id = ?
                AND sleep_date BETWEEN ? AND ?
	),
     second_query
     AS (SELECT DATE,
                Avg(avg_sleep_hours)
                  OVER(
                    ORDER BY DATE ROWS BETWEEN %s preceding AND CURRENT ROW) AS
                   moving_avg
         FROM   first_query)
SELECT *
FROM   second_query
ORDER  BY DATE;`, numDaysPreceding)
	rows, err := slr.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query moving average sleep: %w", err)
	}
	defer rows.Close()

	var movingAverages []models.MovingAverage

	for rows.Next() {
		var movingAvg models.MovingAverage
		var date sql.NullString
		var movingAvgVal sql.NullFloat64

		err := rows.Scan(&date, &movingAvgVal)
		if err != nil {
			return nil, fmt.Errorf("failed to scan moving average sleep row: %w", err)
		}

		if date.Valid {
			d, err := time.Parse(dateLayout, date.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse moving average date: %w", err)
			}
			movingAvg.Date = d
		}

		if movingAvgVal.Valid {
			movingAvg.MovingAvg = movingAvgVal.Float64
		}

		movingAverages = append(movingAverages, movingAvg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating moving average sleep rows: %w", err)
	}

	if movingAverages == nil {
		return make([]models.MovingAverage, 0), nil
	}

	return movingAverages, nil
}

func (slr *SleepLogRepository) StandardDeviation(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	query := `SELECT Stddev_pop(hours_slept) AS std_dev
FROM   sleep_log
WHERE  user_id = ?
       AND sleep_date BETWEEN ? AND ?;`
	rows, err := slr.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("failed to query sleep standard deviation: %w", err)
	}
	defer rows.Close()

	var standardDeviation sql.NullFloat64

	for rows.Next() {
		err := rows.Scan(&standardDeviation)
		if err != nil {
			return 0, fmt.Errorf("failed to scan sleep standard deviation: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating sleep standard deviation rows: %w", err)
	}

	if !standardDeviation.Valid {
		return 0.0, nil
	}

	return standardDeviation.Float64, nil
}

func (slr *SleepLogRepository) SleepQualityTagStat(userID string, startDate time.Time, endDate time.Time) ([]models.TagStat, error) {
	query := `WITH tag_counts AS (
    SELECT 
        sqt.name AS tag_name,
        COUNT(*) AS tag_count
    FROM sleep_log sl
    JOIN sleep_quality_tag sqt 
        ON sl.sleep_quality_tag_id = sqt.sleep_quality_tag_id
    WHERE sl.user_id = ? and sleep_date between ? and ?
    GROUP BY sqt.name
),
tag_percentages AS (
    SELECT
        tag_name,
        tag_count,
        ROUND(tag_count * 100.0 / SUM(tag_count) OVER (), 2) AS percentage
    FROM tag_counts
)
SELECT *
FROM tag_percentages
ORDER BY tag_count ASC;`

	rows, err := slr.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query sleep quality tag stats: %w", err)
	}
	defer rows.Close()

	var sleepTagFrequencies []models.TagStat

	for rows.Next() {
		var sleepTagStat models.TagStat

		err := rows.Scan(
			&sleepTagStat.TagName,
			&sleepTagStat.Count,
			&sleepTagStat.Percentage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sleep quality tag stat row: %w", err)
		}

		sleepTagFrequencies = append(sleepTagFrequencies, sleepTagStat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sleep quality tag stat rows: %w", err)
	}

	if sleepTagFrequencies == nil {
		return make([]models.TagStat, 0), nil
	}

	slices.SortFunc(sleepTagFrequencies, func(a, b models.TagStat) int {
		return int(b.Percentage - a.Percentage)
	})

	return sleepTagFrequencies, nil
}

func (slr *SleepLogRepository) SleepLogs(userID string, startDate time.Time, endDate time.Time) ([]models.SleepLog, error) {
	query := `SELECT sl.sleep_log_id,
       sl.user_id,
       sl.hours_slept,
       sqt.NAME AS sleep_quality_tag_name,
       sl.notes,
       sl.sleep_date,
       sl.created_at,
       sl.updated_at
FROM   sleep_log sl
       JOIN sleep_quality_tag sqt
         ON sl.sleep_quality_tag_id = sqt.sleep_quality_tag_id
WHERE  sl.user_id = ?
       AND sl.sleep_date BETWEEN ? AND ?
ORDER  BY sl.sleep_date;`

	rows, err := slr.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query sleep logs: %w", err)
	}
	defer rows.Close()

	var sleepLogs []models.SleepLog

	for rows.Next() {
		var sleepLog models.SleepLog
		var sleepDateStr, createdAtStr, updatedAtStr string

		err := rows.Scan(
			&sleepLog.SleepLogID,
			&sleepLog.UserID,
			&sleepLog.HoursSlept,
			&sleepLog.SleepQualityTag,
			&sleepLog.Note,
			&sleepDateStr,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sleep log row: %w", err)
		}

		// Parse the date strings into time.Time
		parsedSleepDate, err := time.Parse(dateLayout, sleepDateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sleep_date: %w", err)
		}
		sleepLog.SleepDate = parsedSleepDate

		parsedCreatedAt, err := time.Parse(dateTimeLayout, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		sleepLog.CreatedAt = parsedCreatedAt

		parsedUpdatedAt, err := time.Parse(dateTimeLayout, updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}
		sleepLog.UpdatedAt = parsedUpdatedAt

		sleepLogs = append(sleepLogs, sleepLog)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sleep log rows: %w", err)
	}

	if sleepLogs == nil {
		return make([]models.SleepLog, 0), nil
	}

	return sleepLogs, nil
}

func (slr *SleepLogRepository) DayOfWeekSleepPatterns(userID string, startDate time.Time, endDate time.Time) ([]models.DayOfWeekSleepPattern, error) {
	query := `SELECT
    DAYNAME(sleep_date) AS day_of_week,
    DAYOFWEEK(sleep_date) AS day_number,
    AVG(hours_slept) AS avg_sleep_hours,
    COUNT(*) AS total_entries
FROM sleep_log
WHERE user_id = ?
    AND sleep_date BETWEEN ? AND ?
GROUP BY day_of_week, day_number
ORDER BY day_number;`

	rows, err := slr.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query day-of-week sleep patterns: %w", err)
	}
	defer rows.Close()

	var patterns []models.DayOfWeekSleepPattern

	for rows.Next() {
		var pattern models.DayOfWeekSleepPattern

		err := rows.Scan(
			&pattern.DayOfWeek,
			&pattern.DayNumber,
			&pattern.AvgSleepHours,
			&pattern.TotalEntries,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan day-of-week pattern row: %w", err)
		}

		patterns = append(patterns, pattern)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating day-of-week pattern rows: %w", err)
	}

	if patterns == nil {
		return make([]models.DayOfWeekSleepPattern, 0), nil
	}

	return patterns, nil
}
