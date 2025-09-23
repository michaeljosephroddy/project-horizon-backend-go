package database

import (
	"database/sql"
	"sort"
	"strings"

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

func (jr *JournalRepository) Streaks(userID string, startDate string, endDate string, operator string, moodRating string, moodCategoryID string, targetPercentage string) []models.Streak {
	query := fmt.Sprintf(`WITH qualifying_days AS (
							SELECT Date(je.created_at) AS date,
								   AVG(je.mood_rating) AS avg_rating,
								   COUNT(CASE WHEN mt.mood_category_id = ? THEN 1 END) AS target_category_count,
								   COUNT(mt.mood_tag_id) AS total_count,
								   (COUNT(CASE WHEN mt.mood_category_id = ? THEN 1 END) * 100.0 / COUNT(mt.mood_tag_id)) AS target_percentage
							FROM journal_entry je
							JOIN journal_entry_mood_tag jemt ON je.journal_entry_id = jemt.journal_entry_id
							JOIN mood_tag mt ON jemt.mood_tag_id = mt.mood_tag_id
							WHERE je.user_id = ?
								  AND Date(je.created_at) BETWEEN ? AND ?
							GROUP BY Date(je.created_at)
							HAVING avg_rating %s ? AND target_percentage >= ?
							ORDER BY date
						),
						second_query AS (
							SELECT date, avg_rating,
								   ROW_NUMBER() OVER(ORDER BY date) AS rn
							FROM qualifying_days
						),
						third_query AS (
							SELECT date AS start_date,
								   COUNT(*) AS streak_length,
								   Date_add(date, interval - rn DAY) AS consec_groups
							FROM second_query
							GROUP BY consec_groups
						),
						fourth_query AS (
							SELECT *,
								   Date_add(start_date, interval + streak_length - 1 DAY) AS end_date
								FROM third_query
							)
							SELECT start_date,
								   end_date,
								   streak_length
							FROM fourth_query
							WHERE streak_length >= 2
							ORDER BY start_date;`, operator)

	rows, queryErr := jr.db.Query(query, moodCategoryID, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
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

// TODO need to get mood tags and set it to the list of mood tags on each journalEntry
func (r *JournalRepository) Days(userID string, startDate string, endDate string, operator string, moodRating string, moodCategoryID string, targetPercentage string) []models.Day {
	query := fmt.Sprintf(`SELECT   date,
							created_at,
							journal_entry_id,
							mood_rating,
							note,
							mood_tags,
							mood_tag_ids,
							daily_avg_rating,
							daily_target_count,
							daily_total_count,
							daily_target_percentage
				   FROM (
					   SELECT   Date(je.created_at) AS date,
								je.created_at,
								je.journal_entry_id,
								je.mood_rating,
								je.note,
								GROUP_CONCAT(mt.name ORDER BY mt.name SEPARATOR ', ') AS mood_tags,
								GROUP_CONCAT(mt.mood_tag_id ORDER BY mt.mood_tag_id SEPARATOR ',') AS mood_tag_ids,
								AVG(je.mood_rating) OVER (PARTITION BY Date(je.created_at)) AS daily_avg_rating,
								SUM(CASE WHEN mt.mood_category_id = ? THEN 1 ELSE 0 END) OVER (PARTITION BY Date(je.created_at)) AS daily_target_count,
								COUNT(mt.mood_tag_id) OVER (PARTITION BY Date(je.created_at)) AS daily_total_count,
								(SUM(CASE WHEN mt.mood_category_id = ? THEN 1 ELSE 0 END) OVER (PARTITION BY Date(je.created_at)) * 100.0 / 
								 COUNT(mt.mood_tag_id) OVER (PARTITION BY Date(je.created_at))) AS daily_target_percentage
					   FROM     journal_entry je
					   JOIN     journal_entry_mood_tag jemt
					   ON       je.journal_entry_id = jemt.journal_entry_id
					   JOIN     mood_tag mt
					   ON       jemt.mood_tag_id = mt.mood_tag_id
					   WHERE    je.user_id = ?
					   AND      Date(je.created_at) BETWEEN ? AND ?
					   GROUP BY Date(je.created_at), je.journal_entry_id, je.created_at, je.mood_rating, je.note
				   ) AS daily_data
				   WHERE daily_avg_rating %s ?
				   AND   daily_target_percentage >= ?
				   ORDER BY date, created_at;`, operator)

	rows, queryErr := r.db.Query(query, moodCategoryID, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	resultsByDate := make(map[string]*models.Day)

	for rows.Next() {
		var dateStr string
		var createdAt string
		var journalEntryID int
		var moodRating int
		var note string
		var moodTags string
		var moodTagIDs string
		var dailyAvgRating float64
		var dailyTargetCount int
		var dailyTotalCount int
		var dailyTargetPercentage float64

		err := rows.Scan(
			&dateStr,
			&createdAt,
			&journalEntryID,
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
			panic(err)
		}

		// Create day if it doesn't exist
		if _, exists := resultsByDate[dateStr]; !exists {
			resultsByDate[dateStr] = &models.Day{
				Date:      dateStr,
				AvgRating: dailyAvgRating,
				/* TargetCategoryCount: dailyTargetCount,
				TotalCount:          dailyTotalCount,
				TargetPercentage:    dailyTargetPercentage, */
				JournalEntries: []models.JournalEntry{},
			}
		}

		trimmed := strings.TrimSpace(moodTags)

		parts := strings.Split(trimmed, ",")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		result := strings.Join(parts, ",")
		tags := strings.Split(result, ",")

		// Add journal entry to this day
		entry := models.JournalEntry{
			CreatedAt:      createdAt,
			JournalEntryID: journalEntryID,
			MoodRating:     moodRating,
			Note:           note,
			MoodTags:       tags,
			// MoodTagIds:     moodTagIDs,
		}

		resultsByDate[dateStr].JournalEntries = append(resultsByDate[dateStr].JournalEntries, entry)
	}

	// Convert map to slice
	days := make([]models.Day, 0, len(resultsByDate))
	for _, day := range resultsByDate {
		days = append(days, *day)
	}

	// Sort days by date
	sort.Slice(days, func(i, j int) bool {
		return days[i].Date < days[j].Date
	})

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
