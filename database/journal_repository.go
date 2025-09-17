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

func (jr *JournalRepository) LowDays(userId string, startDate string, endDate string) []models.LowDay {
	q := fmt.Sprintf(`WITH first_query
						 AS (SELECT Date(timestamp)     AS date,
									Avg(overall_rating) AS avg_rating
							 FROM   journal_entry
							 WHERE  user_id = %s 
									AND Date(timestamp) BETWEEN "%s" AND "%s"
							 GROUP  BY Date(timestamp)),
						 second_query
						 AS (SELECT date,
									avg_rating
							 FROM   first_query
							 WHERE  avg_rating <= 4)
					SELECT *
					FROM   second_query;`, userId, startDate, endDate)

	rows, err := jr.db.Query(q)
	if err != nil {
		panic(err)
	}

	var lowDay models.LowDay
	var lowDays []models.LowDay

	for rows.Next() {
		rows.Scan(&lowDay.Date, &lowDay.AvgRating)
		lowDays = append(lowDays, lowDay)
	}

	return lowDays
}

func (jr *JournalRepository) StandardDeviation(userId string, startDate string, endDate string) float32 {

	q := fmt.Sprintf(`SELECT STDDEV_POP(overall_rating) AS std_dev
						FROM   journal_entry
						WHERE  user_id = %s AND DATE(TIMESTAMP) BETWEEN "%s" AND "%s"; 	`, userId, startDate, endDate)

	rows, err := jr.db.Query(q)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var standardDeviation float32
	for rows.Next() {
		rows.Scan(&standardDeviation)
	}

	return standardDeviation
}

func (jr *JournalRepository) MovingAverages(userId string, startDate string, endDate string) []models.MovingAverage {

	q := fmt.Sprintf(`WITH first_query
						 AS (SELECT DATE(TIMESTAMP)     AS DATE,
									AVG(overall_rating) AS daily_avg
							 FROM   journal_entry
							 WHERE  user_id = %s
									AND DATE(TIMESTAMP) BETWEEN "%s" AND "%s"
							 GROUP  BY DATE(TIMESTAMP)),
						 second_query
						 AS (SELECT DATE,
									Avg(daily_avg)
									  OVER(
										ORDER BY DATE ROWS BETWEEN 2 preceding AND CURRENT ROW) AS
									   moving_avg_3day
							 FROM   first_query)
					SELECT *
					FROM   second_query
					ORDER  BY DATE; 	`, userId, startDate, endDate)

	rows, err := jr.db.Query(q)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var movingAverages []models.MovingAverage
	var movingAverage models.MovingAverage

	for rows.Next() {
		rows.Scan(&movingAverage.Date, &movingAverage.ThreeDay)
		movingAverages = append(movingAverages, movingAverage)
	}

	return movingAverages
}

func (jr *JournalRepository) GetAllJournalEntries(userId string) []models.JournalEntry {

	dbQuery := fmt.Sprintf("SELECT * FROM journal_entry WHERE user_id = %s", userId)

	rows, err := jr.db.Query(dbQuery)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var journalEntries []models.JournalEntry
	var journalEntry models.JournalEntry

	for rows.Next() {
		rows.Scan(&journalEntry.JournalEntryId, &journalEntry.UserId,
			&journalEntry.OverallRating, &journalEntry.Note, &journalEntry.Timestamp,
			&journalEntry.CreatedAt, &journalEntry.UpdatedAt)
		journalEntries = append(journalEntries, journalEntry)
	}

	return journalEntries
}

func (jr *JournalRepository) HighDays(userId string, startDate string, endDate string) []models.HighDay {
	query := `WITH first_query
				  AS (SELECT Date(timestamp) AS date,
							 Avg(overall_rating) AS avg_rating
					  FROM journal_entry
					  WHERE user_id = ? 
							AND Date(timestamp) BETWEEN ? AND ?
					  GROUP BY Date(timestamp)),
				  second_query
				  AS (SELECT date, avg_rating
					  FROM first_query
					  WHERE avg_rating >= 6)
				  SELECT date, avg_rating FROM second_query`

	rows, err := jr.db.Query(query, userId, startDate, endDate)
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

func (jr *JournalRepository) MoodTagFrequencies(userId string, startDate string, endDate string) []models.MoodTagFrequency {

	q := fmt.Sprintf(`WITH first_query
						 AS (SELECT journal_entry.journal_entry_id,
									mood_tag_id,
									Date(timestamp)    AS date,
									Count(mood_tag_id) AS mood_tag_id_count
							 FROM   journal_entry
									INNER JOIN journal_entry_mood_tag
											ON journal_entry.journal_entry_id =
											   journal_entry_mood_tag.journal_entry_id
							 WHERE  user_id = %s 
									AND Date(timestamp) BETWEEN "%s" AND "%s"
							 GROUP  BY mood_tag_id),
						 second_query
						 AS (SELECT mood_tag_id_count,
									(mood_tag_id_count / Sum(mood_tag_id_count)
															OVER() ) * 100 AS percentage
							 FROM   first_query)
					SELECT *
					FROM   second_query;`, userId, startDate, endDate)

	rows, err := jr.db.Query(q)
	if err != nil {
		panic(err)
	}

	var moodTagFrequency models.MoodTagFrequency
	var moodTagFrequencies []models.MoodTagFrequency

	for rows.Next() {
		rows.Scan(&moodTagFrequency.Count, &moodTagFrequency.Percentage)
		moodTagFrequencies = append(moodTagFrequencies, moodTagFrequency)
	}

	return moodTagFrequencies
}
