package database

import (
	"database/sql"

	"fmt"
	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type JournalRepository struct {
	db *sql.DB
}

func NewJournalRepository(dbConnection *sql.DB) *JournalRepository {
	return &JournalRepository{
		db: dbConnection,
	}
}

func (jr *JournalRepository) Streaks(userID string, startDate string, endDate string, operator string, val string) []models.Streak {
	query := fmt.Sprintf(`WITH first_query
				 AS (SELECT DATE(created_at) AS DATE,
							AVG(mood_rating) AS avg_rating
					 FROM   journal_entry
					 WHERE  user_id = ?
							AND DATE(created_at) BETWEEN ? AND ? 
					 GROUP  BY DATE(created_at)
					 ORDER  BY DATE(created_at)),
				 second_query
				 AS (SELECT *,
							ROW_NUMBER()
							  OVER(
								ORDER BY DATE) AS rn
					 FROM   first_query
					 WHERE  avg_rating %s ?),
				 third_query
				 AS (SELECT DATE                              AS start_date,
							COUNT(*)                          AS streak_length,
							Date_add(DATE, interval - rn DAY) AS consec_groups
					 FROM   second_query
					 GROUP  BY consec_groups),
				 fourth_query
				 AS (SELECT *,
							Date_add(start_date, interval + streak_length - 1 DAY) AS
							end_date
					 FROM   third_query)
			SELECT start_date,
				   end_date,
				   streak_length
			FROM   fourth_query;`, operator)

	rows, queryErr := jr.db.Query(query, userID, startDate, endDate, val)
	if queryErr != nil {
		panic(queryErr)
	}

	var streak models.Streak
	var streaks []models.Streak

	for rows.Next() {
		scanErr := rows.Scan(&streak.StartDate, &streak.EndDate, &streak.NumDays)
		if scanErr != nil {
			panic(scanErr)
		}
		streaks = append(streaks, streak)
	}

	if streaks == nil {
		return make([]models.Streak, 0)
	}

	return streaks
}

func (jr *JournalRepository) Days(userID string, startDate string, endDate string, operator string, val string) []models.Day {
	query := fmt.Sprintf(`WITH first_query
				 AS (SELECT Date(created_at)     AS date,
							AVG(mood_rating) AS avg_rating
					 FROM   journal_entry
					 WHERE  user_id = ? 
							AND Date(created_at) BETWEEN ? AND ?
					 GROUP  BY Date(created_at)),
				 second_query
				 AS (SELECT date,
							avg_rating
					 FROM   first_query
					 WHERE  avg_rating %s ?)
			SELECT *
			FROM   second_query;`, operator)

	rows, queryErr := jr.db.Query(query, userID, startDate, endDate, val)
	if queryErr != nil {
		panic(queryErr)
	}

	var day models.Day
	var days []models.Day

	for rows.Next() {
		scanErr := rows.Scan(&day.Date, &day.AvgRating)
		if scanErr != nil {
			panic(scanErr)
		}
		days = append(days, day)
	}

	if days == nil {
		return make([]models.Day, 0)
	}

	return days
}

func (jr *JournalRepository) StandardDeviation(userID string, startDate string, endDate string) float32 {
	query := `SELECT STDDEV_POP(mood_rating) AS std_dev
				FROM   journal_entry
				WHERE  user_id = ? AND DATE(created_at) BETWEEN ? AND ?;`

	rows, queryErr := jr.db.Query(query, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}

	defer rows.Close()

	var standardDeviation float32

	for rows.Next() {
		scanErr := rows.Scan(&standardDeviation)
		if scanErr != nil {
			panic(scanErr)
		}
	}

	return standardDeviation
}

func (jr *JournalRepository) MovingAverages(userID string, startDate string, endDate string) []models.MovingAverage {
	query := `WITH first_query
				 AS (SELECT DATE(created_at)     AS DATE,
							AVG(mood_rating) AS daily_avg
					 FROM   journal_entry
					 WHERE  user_id = ? 
							AND DATE(created_at) BETWEEN ? AND ?
					 GROUP  BY DATE(created_at)),
				 second_query
				 AS (SELECT DATE,
							AVG(daily_avg)
							  OVER(
								ORDER BY DATE ROWS BETWEEN 2 preceding AND CURRENT ROW) AS
							   moving_avg_3day
					 FROM   first_query)
			SELECT *
			FROM   second_query
			ORDER  BY DATE;`

	rows, queryErr := jr.db.Query(query, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}

	defer rows.Close()

	var movingAverages []models.MovingAverage
	var movingAverage models.MovingAverage

	for rows.Next() {
		scanErr := rows.Scan(&movingAverage.Date, &movingAverage.ThreeDay)
		if scanErr != nil {
			panic(scanErr)
		}
		movingAverages = append(movingAverages, movingAverage)
	}

	return movingAverages
}

func (jr *JournalRepository) JournalEntries(userID string) []models.JournalEntry {
	query := `SELECT * FROM journal_entry WHERE user_id = ?`

	rows, err := jr.db.Query(query, userID)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var journalEntries []models.JournalEntry
	var journalEntry models.JournalEntry

	for rows.Next() {
		scanErr := rows.Scan(&journalEntry.JournalEntryID, &journalEntry.UserID,
			&journalEntry.MoodRating, &journalEntry.Note, &journalEntry.CreatedAt)
		if scanErr != nil {
			panic(scanErr)
		}
		journalEntries = append(journalEntries, journalEntry)
	}

	return journalEntries
}

func (jr *JournalRepository) MoodTagFrequencies(userID string, startDate string, endDate string) []models.MoodTagFrequency {
	query := `WITH first_query AS (
				SELECT 
					je.journal_entry_id,
					jem.mood_tag_id,
					mt.name, 
					DATE(je.created_at) AS date,
					COUNT(jem.mood_tag_id) AS mood_tag_id_count
				FROM journal_entry je
				INNER JOIN journal_entry_mood_tag jem
					ON je.journal_entry_id = jem.journal_entry_id
				INNER JOIN mood_tag mt
					ON jem.mood_tag_id = mt.mood_tag_id
				WHERE je.user_id = ? 
				  AND DATE(je.created_at) BETWEEN ? AND ?
				GROUP BY jem.mood_tag_id, mt.name, je.journal_entry_id, DATE(je.created_at)
			),
			second_query AS (
				SELECT 
					name,
					SUM(mood_tag_id_count) AS mood_tag_id_count,
					(SUM(mood_tag_id_count) / SUM(SUM(mood_tag_id_count)) OVER()) * 100 AS percentage
				FROM first_query
				GROUP BY mood_tag_id, name
			)
			SELECT *
			FROM second_query;`

	rows, err := jr.db.Query(query, userID, startDate, endDate)
	if err != nil {
		panic(err)
	}

	var moodTagFrequency models.MoodTagFrequency
	var moodTagFrequencies []models.MoodTagFrequency

	for rows.Next() {
		scanErr := rows.Scan(&moodTagFrequency.MoodTag, &moodTagFrequency.Count,
			&moodTagFrequency.Percentage)

		if scanErr != nil {
			panic(scanErr)
		}

		moodTagFrequencies = append(moodTagFrequencies, moodTagFrequency)
	}

	return moodTagFrequencies
}
