package repository

import (
	"database/sql"
	"strings"
	"encoding/json"
	"fmt"
	"math"
	"time"
	"slices"
	"sort"

	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/models"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(dbConnection *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: dbConnection,
	}
}

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
)

func (ar *AnalyticsRepository) Days(userID string, startDate time.Time, endDate time.Time, operator string, moodRating string, moodCategoryID string, targetPercentage string) ([]models.Day, error) {
	// Build the query with the operator
	query := fmt.Sprintf(`WITH mood_data AS
(
         SELECT   Date(ml.created_at) AS date,
                  ml.created_at,
                  ml.mood_log_id,
                  ml.mood_rating,
                  ml.note,
                  group_concat(mt.NAME order BY mt.NAME separator ', ')              AS mood_tags,
                  group_concat(mt.mood_tag_id ORDER BY mt.mood_tag_id separator ',') AS mood_tag_ids,
                  sum(
                  CASE
                           WHEN mt.mood_category_id = ? THEN 1
                           ELSE 0
                  END)                  AS entry_target_count,
                  count(mt.mood_tag_id) AS entry_total_count
         FROM     mood_log ml
         JOIN     mood_log_mood_tag mlmt
         ON       ml.mood_log_id = mlmt.mood_log_id
         JOIN     mood_tag mt
         ON       mlmt.mood_tag_id = mt.mood_tag_id
         WHERE    ml.user_id = ?
         AND      date(ml.created_at) BETWEEN ? AND      ?
         GROUP BY date(ml.created_at),
                  ml.mood_log_id,
                  ml.created_at,
                  ml.mood_rating,
                  ml.note ), daily_stats AS
(
       SELECT date,
              created_at,
              mood_log_id,
              mood_rating,
              note,
              mood_tags,
              mood_tag_ids,
              avg(mood_rating) OVER (partition BY        date)                                                                        AS daily_avg_rating,
              sum(entry_target_count) OVER (partition BY date)                                                                        AS daily_target_count,
              sum(entry_total_count) OVER (partition BY  date)                                                                        AS daily_total_count,
              (sum(entry_target_count) OVER (partition BY date) * 100.0 / NULLIF(sum(entry_total_count) OVER (partition BY date), 0)) AS daily_target_percentage
       FROM   mood_data )
SELECT   date,
         created_at,
         mood_log_id,
         mood_rating,
         note,
         mood_tags,
         mood_tag_ids,
         daily_avg_rating,
         daily_target_count,
         daily_total_count,
         COALESCE(daily_target_percentage, 0) AS daily_target_percentage
FROM     daily_stats
WHERE    daily_avg_rating %s ?
AND      COALESCE(daily_target_percentage, 0) >= ?
ORDER BY date DESC,
         created_at DESC;`, operator)

	// Execute with parameters in correct order
	rows, err := ar.db.Query(query, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to query days: %w", err)
	}
	defer rows.Close()

	resultsByDate := make(map[string]*models.Day)

	for rows.Next() {
		var dateStr string
		var createdAtStr string
		var moodLogID int
		var moodRating int
		var note sql.NullString
		var moodTags sql.NullString
		var moodTagIDs sql.NullString
		var dailyAvgRating float64
		var dailyTargetCount int
		var dailyTotalCount int
		var dailyTargetPercentage float64

		err := rows.Scan(
			&dateStr,
			&createdAtStr,
			&moodLogID,
			&moodRating,
			&note,
			&moodTags,
			&moodTagIDs,
			&dailyAvgRating,
			&dailyTargetCount,
			&dailyTotalCount,
			&dailyTargetPercentage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan day row: %w", err)
		}

		// Parse the date string to time.Time
		parsedDate, err := time.Parse(dateLayout, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		// Parse the createdAt string to time.Time
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			// Try alternative format if RFC3339 fails
			createdAt, err = time.Parse(dateTimeLayout, createdAtStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse created_at: %w", err)
			}
		}

		if _, exists := resultsByDate[dateStr]; !exists {
			resultsByDate[dateStr] = &models.Day{
				Date:           parsedDate,
				DailyAvgRating: dailyAvgRating,
				MoodLogs:       []models.MoodLog{},
			}
		}

		var tags []string
		if moodTags.Valid && moodTags.String != "" {
			mTags := strings.Split(moodTags.String, ", ")
			for _, t := range mTags {
				trimmed := strings.TrimSpace(t)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		noteValue := ""
		if note.Valid {
			noteValue = note.String
		}

		entry := models.MoodLog{
			CreatedAt:  createdAt,
			UserID:     userID,
			MoodLogID:  moodLogID,
			MoodRating: moodRating,
			Note:       noteValue,
			MoodTags:   tags,
		}

		resultsByDate[dateStr].MoodLogs = append(resultsByDate[dateStr].MoodLogs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating day rows: %w", err)
	}

	days := make([]models.Day, 0, len(resultsByDate))
	for _, day := range resultsByDate {
		days = append(days, *day)
	}

	// Sort by date descending (most recent first)
	sort.Slice(days, func(i, j int) bool {
		return days[i].Date.After(days[j].Date)
	})

	return days, nil
}

func (ar *AnalyticsRepository) StandardDeviation(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	query := `SELECT Stddev_pop(mood_rating) AS std_dev
FROM   mood_log
WHERE  user_id = ?
       AND Date(created_at) BETWEEN ? AND ?;`

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("failed to query mood standard deviation: %w", err)
	}
	defer rows.Close()

	var standardDeviation sql.NullFloat64

	for rows.Next() {
		err := rows.Scan(&standardDeviation)
		if err != nil {
			return 0, fmt.Errorf("failed to scan mood standard deviation: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating mood standard deviation rows: %w", err)
	}

	if !standardDeviation.Valid {
		return 0.0, nil
	}

	return standardDeviation.Float64, nil
}

func (ar *AnalyticsRepository) MovingAverages(userID string, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]models.MovingAverage, error) {
	query := fmt.Sprintf(`WITH first_query
     AS (SELECT DATE(created_at) AS DATE,
                Avg(mood_rating) AS daily_avg
         FROM   mood_log
         WHERE  user_id = ?
                AND DATE(created_at) BETWEEN ? AND ?
         GROUP  BY DATE(created_at)),
     second_query
     AS (SELECT DATE,
                Avg(daily_avg)
                  OVER(
                    ORDER BY DATE ROWS BETWEEN %s preceding AND CURRENT ROW) AS
                   moving_avg
         FROM   first_query)
SELECT *
FROM   second_query
ORDER  BY DATE;`, numDaysPreceding)

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query moving averages: %w", err)
	}
	defer rows.Close()

	var movingAverages []models.MovingAverage

	for rows.Next() {
		var dateStr string
		var movingAvg float64

		err := rows.Scan(&dateStr, &movingAvg)
		if err != nil {
			return nil, fmt.Errorf("failed to scan moving average row: %w", err)
		}

		// Parse the date string to time.Time
		parsedDate, err := time.Parse(dateLayout, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse moving average date: %w", err)
		}

		movingAverages = append(movingAverages, models.MovingAverage{
			Date:      parsedDate,
			MovingAvg: movingAvg,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating moving average rows: %w", err)
	}

	if movingAverages == nil {
		return make([]models.MovingAverage, 0), nil
	}

	return movingAverages, nil
}

func (ar *AnalyticsRepository) MoodLogs(userID string, startDate time.Time, endDate time.Time) ([]models.MoodLog, error) {
	query := `SELECT *
FROM   mood_log
WHERE  user_id = ? and DATE(created_at) BETWEEN ? AND ?;`

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query mood logs: %w", err)
	}
	defer rows.Close()

	var moodLogs []models.MoodLog

	for rows.Next() {
		var moodLog models.MoodLog

		err := rows.Scan(
			&moodLog.MoodLogID,
			&moodLog.UserID,
			&moodLog.MoodRating,
			&moodLog.Note,
			&moodLog.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mood log row: %w", err)
		}

		moodLogs = append(moodLogs, moodLog)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating mood log rows: %w", err)
	}

	if moodLogs == nil {
		return make([]models.MoodLog, 0), nil
	}

	return moodLogs, nil
}

func (ar *AnalyticsRepository) MoodTagFrequencies(userID string, startDate time.Time, endDate time.Time) ([]models.TagStat, error) {
	query := `WITH first_query
     AS (SELECT ml.mood_log_id,
                mlmt.mood_tag_id,
                mt.NAME,
                Date(ml.created_at)    AS date,
                Count(mlmt.mood_tag_id) AS mood_tag_id_count
         FROM   mood_log ml
                INNER JOIN mood_log_mood_tag mlmt
                        ON ml.mood_log_id = mlmt.mood_log_id
                INNER JOIN mood_tag mt
                        ON mlmt.mood_tag_id = mt.mood_tag_id
         WHERE  ml.user_id = ?
                AND Date(ml.created_at) BETWEEN ? AND ?
         GROUP  BY mlmt.mood_tag_id,
                   mt.NAME,
                   ml.mood_log_id,
                   Date(ml.created_at)),
     second_query
     AS (SELECT NAME,
                Sum(mood_tag_id_count)                      AS mood_tag_id_count
                ,
                ( Sum(mood_tag_id_count) / Sum(Sum(
                  mood_tag_id_count))
                                             OVER() ) * 100 AS percentage
         FROM   first_query
         GROUP  BY mood_tag_id,
                   NAME)
SELECT *
FROM   second_query;`

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query mood tag frequencies: %w", err)
	}
	defer rows.Close()

	var moodTagFrequencies []models.TagStat

	for rows.Next() {
		var moodTagStat models.TagStat

		err := rows.Scan(
			&moodTagStat.TagName,
			&moodTagStat.Count,
			&moodTagStat.Percentage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mood tag frequency row: %w", err)
		}

		moodTagFrequencies = append(moodTagFrequencies, moodTagStat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating mood tag frequency rows: %w", err)
	}

	slices.SortFunc(moodTagFrequencies, func(a, b models.TagStat) int {
		return int(b.Percentage - a.Percentage)
	})

	if moodTagFrequencies == nil {
		return make([]models.TagStat, 0), nil
	}

	return moodTagFrequencies, nil
}

func (ar *AnalyticsRepository) AvgMoodRating(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	query := `WITH first_query
     AS (SELECT Date(created_at) AS date,
                AVG(mood_rating) AS daily_avg_rating
         FROM   mood_log
         WHERE  user_id = ? 
                AND Date(created_at) BETWEEN ? AND ? 
         GROUP  BY Date(created_at)),
     second_query
     AS (SELECT AVG(daily_avg_rating) AS period_mood_rating_avg
         FROM   first_query)
SELECT period_mood_rating_avg
FROM   second_query;`

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("failed to query average mood rating: %w", err)
	}
	defer rows.Close()

	var avgMoodRatingPeriod sql.NullFloat64

	for rows.Next() {
		err := rows.Scan(&avgMoodRatingPeriod)
		if err != nil {
			return 0, fmt.Errorf("failed to scan average mood rating: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating average mood rating rows: %w", err)
	}

	if !avgMoodRatingPeriod.Valid {
		return 0.0, nil
	}

	return avgMoodRatingPeriod.Float64, nil
}
func (ar *AnalyticsRepository) AvgSleepHours(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	query := `SELECT Avg(hours_slept) AS avg_sleep_hours
FROM   sleep_log
WHERE  user_id = ?
       AND sleep_date BETWEEN ? AND ?;`

	rows, err := ar.db.Query(query, userID, startDate, endDate)
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

func (ar *AnalyticsRepository) MovingAvgSleep(userID string, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]models.MovingAverage, error) {
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
	rows, err := ar.db.Query(query, userID, startDate, endDate)
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

func (ar *AnalyticsRepository) SleepStandardDeviation(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	query := `SELECT Stddev_pop(hours_slept) AS std_dev
FROM   sleep_log
WHERE  user_id = ?
       AND sleep_date BETWEEN ? AND ?;`
	rows, err := ar.db.Query(query, userID, startDate, endDate)
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

func (ar *AnalyticsRepository) SleepQualityTagStat(userID string, startDate time.Time, endDate time.Time) ([]models.TagStat, error) {
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

	rows, err := ar.db.Query(query, userID, startDate, endDate)
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

func (ar *AnalyticsRepository) SleepLogs(userID string, startDate time.Time, endDate time.Time) ([]models.SleepLog, error) {
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

	rows, err := ar.db.Query(query, userID, startDate, endDate)
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

func (ar *AnalyticsRepository) DayOfWeekSleepPatterns(userID string, startDate time.Time, endDate time.Time) ([]models.DayOfWeekSleepPattern, error) {
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

	rows, err := ar.db.Query(query, userID, startDate, endDate)
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
func (ar *AnalyticsRepository) MedicationLogs(userID string, startDate, endDate time.Time) ([]models.MedicationLog, error) {
	query := `SELECT 
		ml.medication_log_id,
		ml.user_id,
		ml.taken_at,
		ml.notes AS log_notes,
		ml.created_at,
		ml.updated_at,
		JSON_ARRAYAGG(
			JSON_OBJECT(
				'medication_id', m.medication_id,
				'name', m.name,
				'dosage', mli.dosage
			)
		) AS medications
	FROM medication_log ml
	JOIN medication_log_item mli ON ml.medication_log_id = mli.medication_log_id
	JOIN medication m ON m.medication_id = mli.medication_id
	WHERE ml.user_id = ?
		AND DATE(ml.taken_at) BETWEEN ? AND ?
	GROUP BY 
		ml.medication_log_id,
		ml.user_id,
		ml.taken_at,
		ml.notes,
		ml.created_at,
		ml.updated_at
	ORDER BY ml.taken_at DESC`

	rows, err := ar.db.Query(query, userID, startDate, endDate)
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
func (ar *AnalyticsRepository) OverviewStats(userID string, startDate, endDate time.Time) (float64, error) {
	query := `SELECT 
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
			AND DATE(ml.taken_at) BETWEEN ? AND ?`

	var daysWithLogs, totalDays int
	var avgMedsPerLog sql.NullFloat64

	err := ar.db.QueryRow(query, endDate, startDate, userID, startDate, endDate).Scan(
		&daysWithLogs, &totalDays, &avgMedsPerLog,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get overview stats: %w", err)
	}

	adherenceRate := 0.0
	if totalDays > 0 {
		adherenceRate = (float64(daysWithLogs) / float64(totalDays)) * 100
	}

	return adherenceRate, nil
}

// MedicationDetailedStats returns comprehensive stats for each medication
func (ar *AnalyticsRepository) MedicationDetailedStats(userID string, startDate, endDate time.Time) ([]models.MedicationStats, error) {
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

	rows, err := ar.db.Query(query, endDate, startDate, userID, startDate, endDate)
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
		timingStats, err := ar.getTimingStats(userID, medID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		stat.AvgTakenAtTime = timingStats.AvgTime
		stat.TimingStdDevMinutes = timingStats.StdDevMinutes
		stat.TimingDescription = timingStats.Description
		stat.EarliestTime = timingStats.EarliestTime
		stat.LatestTime = timingStats.LatestTime

		// Get streaks
		longestStreak, currentStreak, err := ar.getStreaks(userID, medID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		stat.LongestStreak = longestStreak
		stat.CurrentStreak = currentStreak

		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

func (ar *AnalyticsRepository) getTimingStats(userID string, medID int, startDate, endDate time.Time) (models.TimingStats, error) {
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

	err := ar.db.QueryRow(query, userID, medID, startDate, endDate).Scan(
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

func (ar *AnalyticsRepository) getStreaks(userID string, medID int, startDate, endDate time.Time) (int, int, error) {
	query := `SELECT DISTINCT DATE(ml.taken_at) as log_date
		FROM medication_log ml
		JOIN medication_log_item mli ON ml.medication_log_id = mli.medication_log_id
		WHERE ml.user_id = ?
			AND mli.medication_id = ?
			AND DATE(ml.taken_at) BETWEEN ? AND ?
		ORDER BY log_date`

	rows, err := ar.db.Query(query, userID, medID, startDate, endDate)
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
