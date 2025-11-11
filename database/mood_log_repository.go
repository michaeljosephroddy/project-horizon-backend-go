package database

import (
	"database/sql"
	"slices"
	"sort"
	"strings"
	"time"

	"fmt"

	"github.com/michaeljosephroddy/project-horizon-backend-go/analytics/models"
)

type MoodLogRepository struct {
	db *sql.DB
}

func NewMoodLogRepository(dbConnection *sql.DB) *MoodLogRepository {
	return &MoodLogRepository{
		db: dbConnection,
	}
}

// NOT in use atm
/* func (mlr *MoodLogRepository) Streaks(userID string, startDate time.Time, endDate time.Time, operator string, moodRating string, moodCategoryID string, targetPercentage string) ([]models.Streak, error) {
	query := fmt.Sprintf(streaksQuery, operator)
	rows, err := mlr.db.Query(query, moodCategoryID, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to query streaks: %w", err)
	}
	defer rows.Close()

	var streaks []models.Streak

	for rows.Next() {
		var streak models.Streak
		var startDateStr, endDateStr string

		err := rows.Scan(
			&startDateStr,
			&endDateStr,
			&streak.NumDays,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan streak row: %w", err)
		}

		// Parse the date strings into time.Time
		parsedStartDate, err := time.Parse(dateLayout, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse streak start date: %w", err)
		}

		parsedEndDate, err := time.Parse(dateLayout, endDateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse streak end date: %w", err)
		}

		streak.StartDate = parsedStartDate
		streak.EndDate = parsedEndDate

		streaks = append(streaks, streak)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating streak rows: %w", err)
	}

	for i := 0; i < len(streaks); i++ {
		streakDays, err := mlr.Days(userID, streaks[i].StartDate, streaks[i].EndDate, operator, moodRating, moodCategoryID, targetPercentage)
		if err != nil {
			return nil, fmt.Errorf("failed to get days for streak: %w", err)
		}
		streaks[i].Days = append(streaks[i].Days, streakDays...)
	}

	if streaks == nil {
		return make([]models.Streak, 0), nil
	}

	return streaks, nil
} */

func (mlr *MoodLogRepository) Days(userID string, startDate time.Time, endDate time.Time, operator string, moodRating string, moodCategoryID string, targetPercentage string) ([]models.Day, error) {
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
	rows, err := mlr.db.Query(query, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
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

func (mlr *MoodLogRepository) StandardDeviation(userID string, startDate time.Time, endDate time.Time) (float64, error) {
	query := `SELECT Stddev_pop(mood_rating) AS std_dev
FROM   mood_log
WHERE  user_id = ?
       AND Date(created_at) BETWEEN ? AND ?;`

	rows, err := mlr.db.Query(query, userID, startDate, endDate)
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

func (mlr *MoodLogRepository) MovingAverages(userID string, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]models.MovingAverage, error) {
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

	rows, err := mlr.db.Query(query, userID, startDate, endDate)
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

func (mlr *MoodLogRepository) MoodLogs(userID string, startDate time.Time, endDate time.Time) ([]models.MoodLog, error) {
	query := `SELECT *
FROM   mood_log
WHERE  user_id = ? and DATE(created_at) BETWEEN ? AND ?;`

	rows, err := mlr.db.Query(query, userID, startDate, endDate)
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

func (mlr *MoodLogRepository) MoodTagFrequencies(userID string, startDate time.Time, endDate time.Time) ([]models.TagStat, error) {
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

	rows, err := mlr.db.Query(query, userID, startDate, endDate)
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

func (mlr *MoodLogRepository) AvgMoodRating(userID string, startDate time.Time, endDate time.Time) (float64, error) {
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

	rows, err := mlr.db.Query(query, userID, startDate, endDate)
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
