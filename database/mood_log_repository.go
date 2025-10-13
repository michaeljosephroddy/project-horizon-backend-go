package database

import (
	"database/sql"
	"slices"
	"sort"
	"strings"

	"fmt"
	"log"

	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type MoodLogRepository struct {
	db *sql.DB
}

func NewMoodLogRepository(dbConnection *sql.DB) *MoodLogRepository {
	return &MoodLogRepository{
		db: dbConnection,
	}
}

func (mlr *MoodLogRepository) Streaks(userID string, startDate string, endDate string, operator string, moodRating string, moodCategoryID string, targetPercentage string) []models.Streak {

	query := fmt.Sprintf(streaksQuery, operator)
	rows, queryErr := mlr.db.Query(query, moodCategoryID, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
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
		streakDays := mlr.Days(userID, streaks[i].StartDate, streaks[i].EndDate, operator, moodRating, moodCategoryID, targetPercentage)
		streaks[i].Days = append(streaks[i].Days, streakDays...)
	}

	if streaks == nil {
		return make([]models.Streak, 0)
	}

	return streaks
}

func (mlr *MoodLogRepository) Days(userID string, startDate string, endDate string, operator string, moodRating string, moodCategoryID string, targetPercentage string) []models.Day {
	// Build the query with the operator
	query := fmt.Sprintf(daysQuery, operator)

	// Execute with parameters in correct order
	rows, queryErr := mlr.db.Query(query, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
	if queryErr != nil {
		log.Printf("Query error: %v\nQuery: %s\nParams: %v, %v, %v, %v, %v, %v",
			queryErr, query, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
		panic(queryErr)
	}
	defer rows.Close()

	resultsByDate := make(map[string]*models.Day)

	for rows.Next() {
		var dateStr string
		var createdAt string
		var moodLogID int
		var moodRating int
		var note sql.NullString
		var moodTags sql.NullString
		var moodTagIDs sql.NullString
		var dailyAvgRating float64
		var dailyTargetCount int
		var dailyTotalCount int
		var dailyTargetPercentage float64

		scanErr := rows.Scan(
			&dateStr,
			&createdAt,
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
		if scanErr != nil {
			panic(scanErr)
		}

		if _, exists := resultsByDate[dateStr]; !exists {
			resultsByDate[dateStr] = &models.Day{
				Date:           dateStr,
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

	days := make([]models.Day, 0, len(resultsByDate))
	for _, day := range resultsByDate {
		days = append(days, *day)
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i].Date > days[j].Date
	})

	return days
}

func (mlr *MoodLogRepository) StandardDeviation(userID string, startDate string, endDate string) float64 {

	rows, queryErr := mlr.db.Query(stdDevQuery, userID, startDate, endDate)
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

func (mlr *MoodLogRepository) MovingAverages(userID string, startDate string, endDate string, numDaysPreceding string) []models.MovingAverage {

	query := fmt.Sprintf(moodMovingAvgQuery, numDaysPreceding)
	rows, queryErr := mlr.db.Query(query, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var movingAverages []models.MovingAverage
	var movingAverage models.MovingAverage

	for rows.Next() {
		scanErr := rows.Scan(
			&movingAverage.Date,
			&movingAverage.MovingAvg,
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

func (mlr *MoodLogRepository) MoodLogs(userID string, startDate string, endDate string) []models.MoodLog {

	rows, err := mlr.db.Query(journalEntriesQuery, userID, startDate, endDate)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var moodLogs []models.MoodLog
	var moodLog models.MoodLog

	for rows.Next() {
		scanErr := rows.Scan(
			&moodLog.MoodLogID,
			&moodLog.UserID,
			&moodLog.MoodRating,
			&moodLog.Note,
			&moodLog.CreatedAt,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		moodLogs = append(moodLogs, moodLog)
	}

	if moodLogs == nil {
		return make([]models.MoodLog, 0)
	}

	return moodLogs
}

func (mlr *MoodLogRepository) MoodTagFrequencies(userID string, startDate string, endDate string) []models.TagStat {

	rows, queryErr := mlr.db.Query(moodTagFrequenciesQuery, userID, startDate, endDate)
	if queryErr != nil {
		panic(queryErr)
	}
	defer rows.Close()

	var moodTagStat models.TagStat
	var moodTagFrequencies []models.TagStat

	for rows.Next() {
		scanErr := rows.Scan(
			&moodTagStat.TagName,
			&moodTagStat.Count,
			&moodTagStat.Percentage,
		)
		if scanErr != nil {
			panic(scanErr)
		}

		moodTagFrequencies = append(moodTagFrequencies, moodTagStat)
	}

	slices.SortFunc(moodTagFrequencies, func(a, b models.TagStat) int {
		return int(b.Percentage - a.Percentage)
	})

	if moodTagFrequencies == nil {
		return make([]models.TagStat, 0)
	}

	return moodTagFrequencies
}

func (mlr *MoodLogRepository) AvgMoodRating(userID string, startDate string, endDate string) float64 {

	rows, queryErr := mlr.db.Query(AvgMoodRatingQuery, userID, startDate, endDate)
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
