package database

import (
	"database/sql"
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

func (jr *JournalRepository) LowDays(userID string, startDate string, endDate string) []models.LowDay {
	query := `WITH first_query
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
					 WHERE  avg_rating <= 4)
			SELECT *
			FROM   second_query;`

	rows, queryErr := jr.db.Query(query, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}

	var lowDay models.LowDay
	var lowDays []models.LowDay

	for rows.Next() {
		scanErr := rows.Scan(&lowDay.Date, &lowDay.AvgRating)
		if scanErr != nil {
			panic(scanErr)
		}
		lowDays = append(lowDays, lowDay)
	}

	return lowDays
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

func (jr *JournalRepository) HighDays(userID string, startDate string, endDate string) []models.HighDay {
	query := `WITH first_query
				  AS (SELECT Date(created_at) AS date,
							 AVG(mood_rating) AS avg_rating
					  FROM journal_entry
					  WHERE user_id = ? 
							AND Date(created_at) BETWEEN ? AND ?
					  GROUP BY Date(created_at)),
				  second_query
				  AS (SELECT date, avg_rating
					  FROM first_query
					  WHERE avg_rating >= 6)
				  SELECT date, avg_rating FROM second_query`

	rows, err := jr.db.Query(query, userID, startDate, endDate)
	if err != nil {
		panic(err)
	}

	var highDay models.HighDay
	var highDays []models.HighDay

	for rows.Next() {
		rows.Scan(&highDay.Date, &highDay.AvgRating)
		highDays = append(highDays, highDay)
	}

	return highDays
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
					mood_tag_id,
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
		scanErr := rows.Scan(&moodTagFrequency.MoodTagID, &moodTagFrequency.Name, &moodTagFrequency.Count, &moodTagFrequency.Percentage)
		if scanErr != nil {
			panic(scanErr)
		}
		moodTagFrequencies = append(moodTagFrequencies, moodTagFrequency)
	}

	return moodTagFrequencies
}
