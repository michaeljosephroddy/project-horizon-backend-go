package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/analytics/model"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) IAnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
)

var home = os.Getenv("HOME")
var queriesDir = "/repos/project-horizon-backend-go/internal/domain/analytics/repository/queries/"

func (ar *AnalyticsRepository) Days(userID int, startDate time.Time, endDate time.Time, operator string, moodRating string, moodCategoryID string, targetPercentage string) ([]model.Day, error) {

	fileName := "days.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.Day, 0, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := fmt.Sprintf(string(content), operator)

	// Execute with parameters in correct order
	rows, err := ar.db.Query(query, moodCategoryID, userID, startDate, endDate, moodRating, targetPercentage)
	if err != nil {
		return nil, fmt.Errorf("failed to query days: %w", err)
	}
	defer rows.Close()

	resultsByDate := make(map[string]*model.Day)

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
			resultsByDate[dateStr] = &model.Day{
				Date:           parsedDate,
				DailyAvgRating: dailyAvgRating,
				MoodLogs:       []model.MoodLog{},
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

		entry := model.MoodLog{
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

	days := make([]model.Day, 0, len(resultsByDate))
	for _, day := range resultsByDate {
		days = append(days, *day)
	}

	// Sort by date descending (most recent first)
	sort.Slice(days, func(i, j int) bool {
		return days[i].Date.After(days[j].Date)
	})

	return days, nil
}

func (ar *AnalyticsRepository) StandardDeviationMood(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	fileName := "standard_deviation.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

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

func (ar *AnalyticsRepository) MovingAverages(userID int, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error) {
	fileName := "moving_averages.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.MovingAverage, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := fmt.Sprintf(string(content), numDaysPreceding)

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query moving averages: %w", err)
	}
	defer rows.Close()

	var movingAverages []model.MovingAverage

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

		movingAverages = append(movingAverages, model.MovingAverage{
			Date:      parsedDate,
			MovingAvg: movingAvg,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating moving average rows: %w", err)
	}

	if movingAverages == nil {
		return make([]model.MovingAverage, 0), nil
	}

	return movingAverages, nil
}

func (ar *AnalyticsRepository) MoodLogs(userID int, startDate time.Time, endDate time.Time) ([]model.MoodLog, error) {
	// Load query file
	fileName := "mood_logs.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.MoodLog, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query mood logs: %w", err)
	}
	defer rows.Close()

	var moodLogs []model.MoodLog

	for rows.Next() {
		var moodLog model.MoodLog

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
		return make([]model.MoodLog, 0), nil
	}

	return moodLogs, nil
}

func (ar *AnalyticsRepository) MoodTagStats(userID int, startDate time.Time, endDate time.Time) ([]model.TagStat, error) {
	// Load query file
	fileName := "mood_tag_stats.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.TagStat, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query mood tag frequencies: %w", err)
	}
	defer rows.Close()

	var moodTagFrequencies []model.TagStat

	for rows.Next() {
		var moodTagStat model.TagStat

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

	slices.SortFunc(moodTagFrequencies, func(a, b model.TagStat) int {
		return int(b.Percentage - a.Percentage)
	})

	if moodTagFrequencies == nil {
		return make([]model.TagStat, 0), nil
	}

	return moodTagFrequencies, nil
}

func (ar *AnalyticsRepository) AvgMoodRating(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	// Load query file
	fileName := "avg_mood_rating.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

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
func (ar *AnalyticsRepository) AvgSleepHours(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	// Load query file
	fileName := "avg_sleep_hours.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

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

func (ar *AnalyticsRepository) MovingAvgSleep(userID int, startDate time.Time, endDate time.Time, numDaysPreceding string) ([]model.MovingAverage, error) {
	fileName := "moving_avg_sleep.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.MovingAverage, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := fmt.Sprintf(string(content), numDaysPreceding)
	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query moving average sleep: %w", err)
	}
	defer rows.Close()

	var movingAverages []model.MovingAverage

	for rows.Next() {
		var movingAvg model.MovingAverage
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
		return make([]model.MovingAverage, 0), nil
	}

	return movingAverages, nil
}

func (ar *AnalyticsRepository) StandardDeviationSleep(userID int, startDate time.Time, endDate time.Time) (float64, error) {
	// Load query file
	fileName := "standard_deviation_sleep.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

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

func (ar *AnalyticsRepository) SleepQualityTagStats(userID int, startDate time.Time, endDate time.Time) ([]model.TagStat, error) {
	// Load query file
	fileName := "sleep_quality_tag_stats.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.TagStat, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query sleep quality tag stats: %w", err)
	}
	defer rows.Close()

	var sleepTagFrequencies []model.TagStat

	for rows.Next() {
		var sleepTagStat model.TagStat

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
		return make([]model.TagStat, 0), nil
	}

	slices.SortFunc(sleepTagFrequencies, func(a, b model.TagStat) int {
		return int(b.Percentage - a.Percentage)
	})

	return sleepTagFrequencies, nil
}

func (ar *AnalyticsRepository) SleepLogs(userID int, startDate time.Time, endDate time.Time) ([]model.SleepLog, error) {
	// Load query file
	fileName := "sleep_logs.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.SleepLog, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query sleep logs: %w", err)
	}
	defer rows.Close()

	var sleepLogs []model.SleepLog

	for rows.Next() {
		var sleepLog model.SleepLog
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
		return make([]model.SleepLog, 0), nil
	}

	return sleepLogs, nil
}

func (ar *AnalyticsRepository) DayOfWeekSleepPatterns(userID int, startDate time.Time, endDate time.Time) ([]model.DayOfWeekSleepPattern, error) {
	// Load query file
	fileName := "day_of_week_sleep_pattern.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.DayOfWeekSleepPattern, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query day-of-week sleep patterns: %w", err)
	}
	defer rows.Close()

	var patterns []model.DayOfWeekSleepPattern

	for rows.Next() {
		var pattern model.DayOfWeekSleepPattern

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
		return make([]model.DayOfWeekSleepPattern, 0), nil
	}

	return patterns, nil
}
func (ar *AnalyticsRepository) MedicationLogs(userID int, startDate, endDate time.Time) ([]model.MedicationLog, error) {
	// Load query file
	fileName := "day_of_week_sleep_pattern.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.MedicationLog, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	rows, err := ar.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query medication logs: %w", err)
	}
	defer rows.Close()

	var medicationLogs []model.MedicationLog

	for rows.Next() {
		var medicationLog model.MedicationLog
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
		var meds []model.Medication
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
func (ar *AnalyticsRepository) MedicationOverviewStats(userID int, startDate, endDate time.Time) (float64, error) {
	// Load query file
	fileName := "medication_overview_stats.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	var daysWithLogs, totalDays int
	var avgMedsPerLog sql.NullFloat64

	queryErr := ar.db.QueryRow(query, endDate, startDate, userID, startDate, endDate).Scan(
		&daysWithLogs, &totalDays, &avgMedsPerLog,
	)
	if queryErr != nil {
		return 0, fmt.Errorf("failed to get overview stats: %w", err)
	}

	adherenceRate := 0.0
	if totalDays > 0 {
		adherenceRate = (float64(daysWithLogs) / float64(totalDays)) * 100
	}

	return adherenceRate, nil
}

// MedicationDetailedStats returns comprehensive stats for each medication
func (ar *AnalyticsRepository) MedicationDetailedStats(userID int, startDate, endDate time.Time) ([]model.MedicationStats, error) {
	// Load query file
	fileName := "medication_detailed_stats.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.MedicationStats, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	rows, err := ar.db.Query(query, endDate, startDate, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication stats: %w", err)
	}
	defer rows.Close()

	var stats []model.MedicationStats

	for rows.Next() {
		var medID, totalDoses, daysActive, totalDays int
		var name string

		if err := rows.Scan(&medID, &name, &totalDoses, &daysActive, &totalDays); err != nil {
			return nil, fmt.Errorf("failed to scan medication stat: %w", err)
		}

		stat := model.MedicationStats{
			MedicationID:   medID,
			Name:           name,
			TotalDoses:     totalDoses,
			DaysActive:     daysActive,
			AvgDosesPerDay: float64(totalDoses) / float64(totalDays),
		}

		// Get timing stats
		timingStats, err := ar.MedicationTimingStats(userID, medID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		stat.AvgTakenAtTime = timingStats.AvgTime
		stat.TimingStdDevMinutes = timingStats.StdDevMinutes
		stat.TimingDescription = timingStats.Description
		stat.EarliestTime = timingStats.EarliestTime
		stat.LatestTime = timingStats.LatestTime

		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

func (ar *AnalyticsRepository) MedicationTimingStats(userID int, medID int, startDate, endDate time.Time) (model.TimingStats, error) {
	// Load query file
	fileName := "medication_timing_stats.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return model.TimingStats{}, fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)	

	var avgTime, earliestTime, latestTime string
	var stdDevMinutes sql.NullFloat64

	queryErr := ar.db.QueryRow(query, userID, medID, startDate, endDate).Scan(
		&avgTime, &stdDevMinutes, &earliestTime, &latestTime,
	)
	if queryErr != nil {
		return model.TimingStats{}, fmt.Errorf("failed to get timing stats: %w", err)
	}

	stdDev := 0.0
	if stdDevMinutes.Valid {
		stdDev = stdDevMinutes.Float64
	}

	// Format description like "8:47 AM ± 45 minutes"
	description := formatTimingDescription(avgTime, stdDev)

	return model.TimingStats{
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
