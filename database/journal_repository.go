package database

import (
	"database/sql"
	"slices"
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

	query := fmt.Sprintf(streaksQueryStr, operator)
	rows, queryErr := jr.db.Query(query, moodCategoryID, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var streak models.Streak
	var streaks []models.Streak

	for rows.Next() {
		scanErr := rows.Scan(
			&streak.StartDate,
			&streak.EndDate,
			&streak.NumDays,
		)
		if scanErr != nil {
			panic(scanErr)
		}
		streaks = append(streaks, streak)
	}

	for i := 0; i < len(streaks); i++ {
		streakDays := jr.Days(userID, streaks[i].StartDate, streaks[i].EndDate, operator, moodRating, moodCategoryID, targetPercentage)
		streaks[i].Days = append(streaks[i].Days, streakDays...)
	}

	if streaks == nil {
		return make([]models.Streak, 0)
	}

	return streaks
}

func (jr *JournalRepository) Days(userID string, startDate string, endDate string, operator string, moodRating string, moodCategoryID string, targetPercentage string) []models.Day {

	query := fmt.Sprintf(daysQueryStr, operator)

	rows, queryErr := jr.db.Query(query, moodCategoryID, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
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

		scanErr := rows.Scan(
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
		if scanErr != nil {
			panic(scanErr)
		}

		// Create day if it doesn't exist
		if _, exists := resultsByDate[dateStr]; !exists {
			resultsByDate[dateStr] = &models.Day{
				Date:           dateStr,
				DailyAvgRating: dailyAvgRating,
				JournalEntries: []models.JournalEntry{},
			}
		}

		var tags []string
		mTags := strings.Split(moodTags, ",")
		for _, t := range mTags {
			trimmed := strings.TrimSpace(t)
			tags = append(tags, trimmed)
		}

		// Add journal entry to this day
		entry := models.JournalEntry{
			CreatedAt:      createdAt,
			UserID:         userID,
			JournalEntryID: journalEntryID,
			MoodRating:     moodRating,
			Note:           note,
			MoodTags:       tags,
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

	// get daily mood tag frequencies
	for i := 0; i < len(days); i++ {
		dailyMoodTagFrequencies := jr.MoodTagFrequencies(userID, days[i].Date, days[i].Date)
		slices.SortFunc(dailyMoodTagFrequencies, func(a, b models.MoodTagFrequency) int {
			if a.Percentage > b.Percentage {
				return -1
			} else if a.Percentage < b.Percentage {
				return 1
			} else {
				return 0
			}
		})

		days[i].MoodTagFrequencies = append(days[i].MoodTagFrequencies, dailyMoodTagFrequencies...)
	}

	if days == nil {
		return make([]models.Day, 0)
	}

	return days
}

func (jr *JournalRepository) StandardDeviation(userID string, startDate string, endDate string) float64 {

	rows, queryErr := jr.db.Query(stdDevQueryStr, userID, startDate, endDate)
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

func (jr *JournalRepository) MovingAverages(userID string, startDate string, endDate string) []models.MovingAverage {

	rows, queryErr := jr.db.Query(movingAvgQueryStr, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var movingAverages []models.MovingAverage
	var movingAverage models.MovingAverage

	for rows.Next() {
		scanErr := rows.Scan(
			&movingAverage.Date,
			&movingAverage.ThreeDay,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		movingAverages = append(movingAverages, movingAverage)
	}

	if movingAverages == nil {
		return make([]models.MovingAverage, 0)
	}

	return movingAverages
}

func (jr *JournalRepository) JournalEntries(userID string) []models.JournalEntry {

	rows, err := jr.db.Query(journalEntriesQueryStr, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var journalEntries []models.JournalEntry
	var journalEntry models.JournalEntry

	for rows.Next() {
		scanErr := rows.Scan(
			&journalEntry.JournalEntryID,
			&journalEntry.UserID,
			&journalEntry.MoodRating,
			&journalEntry.Note,
			&journalEntry.CreatedAt,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		journalEntries = append(journalEntries, journalEntry)
	}

	if journalEntries == nil {
		return make([]models.JournalEntry, 0)
	}

	return journalEntries
}

func (jr *JournalRepository) MoodTagFrequencies(userID string, startDate string, endDate string) []models.MoodTagFrequency {

	rows, queryErr := jr.db.Query(moodTagFrequenciesQueryStr, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var moodTagFrequency models.MoodTagFrequency
	var moodTagFrequencies []models.MoodTagFrequency

	for rows.Next() {
		scanErr := rows.Scan(
			&moodTagFrequency.MoodTag,
			&moodTagFrequency.Count,
			&moodTagFrequency.Percentage,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		moodTagFrequencies = append(moodTagFrequencies, moodTagFrequency)
	}

	if moodTagFrequencies == nil {
		return make([]models.MoodTagFrequency, 0)
	}

	return moodTagFrequencies
}

func (jr *JournalRepository) PeriodAvgMoodRating(userID string, startDate string, endDate string) float64 {

	rows, queryErr := jr.db.Query(AvgMoodRatingPeriodQueryStr, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var avgMoodRatingPeriod sql.NullFloat64
	for rows.Next() {
		scanErr := rows.Scan(&avgMoodRatingPeriod)
		if scanErr != nil {
			panic(scanErr)
		}
	}

	if !avgMoodRatingPeriod.Valid {
		return 0.0
	}

	return avgMoodRatingPeriod.Float64
}
