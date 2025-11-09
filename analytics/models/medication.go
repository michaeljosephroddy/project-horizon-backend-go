package models

import "time"

// Medication
type Medication struct {
	MedicationID int    `json:"medicationId"`
	Name         string `json:"name"`
	Dosage       string `json:"dosage"`
}

// Medication Log
type MedicationLog struct {
	MedicationLogID int          `json:"medicationLogId"`
	UserID          int          `json:"userId"`
	TakenAt         time.Time    `json:"takenAt"` // date stored as string
	Note            string       `json:"note"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	Medications     []Medication `json:"medications"`
}

// MedicationMetric represents comprehensive medication analytics
type MedicationMetric struct {
	UserID            string                 `json:"userID"`
	Granularity       string                 `json:"granularity"`
	StartDate         time.Time              `json:"startDate"`
	EndDate           time.Time              `json:"endDate"`
	// Overview Stats
	TotalLogs         int                    `json:"totalLogs"`
	AdherenceRate     float64                `json:"adherenceRate"`    // % of days with logs
	AvgLogsPerDay     float64                `json:"avgLogsPerDay"`
	// Per-Medication Details
	MedicationStats   []MedicationStats      `json:"medicationStats"`
	// Polypharmacy
	AvgMedsPerLog     float64                `json:"avgMedsPerLog"`
	// Comparison to previous period
	MedicationDiffs   MedicationDiff         `json:"medicationDiffs"`
}

// MedicationStats represents detailed analytics for a specific medication
type MedicationStats struct {
	MedicationID          int                `json:"medicationId"`
	Name                  string             `json:"name"`
	TotalDoses            int                `json:"totalDoses"`
	DaysActive            int                `json:"daysActive"`
	AvgDosesPerDay        float64            `json:"avgDosesPerDay"`
	// Timing Analysis (raw, no windows)
	AvgTakenAtTime        string             `json:"avgTakenAtTime"`        // "08:30:15"
	TimingStdDevMinutes   float64            `json:"timingStdDevMinutes"`   // e.g., 45.3 minutes
	TimingDescription     string             `json:"timingDescription"`     // "8:47 AM ± 45 minutes"
	EarliestTime          string             `json:"earliestTime"`          // "06:15:00"
	LatestTime            string             `json:"latestTime"`            // "22:30:00"
	// Streaks
	LongestStreak         int                `json:"longestStreak"`         // Consecutive days taken
	CurrentStreak         int                `json:"currentStreak"`
}

// MedicationDiff compares current period to previous period
type MedicationDiff struct {
	TotalLogs             MetricChange              `json:"totalLogs"`
	AdherenceRate         MetricChange              `json:"adherenceRate"`
	AvgLogsPerDay         MetricChange              `json:"avgLogsPerDay"`
	AvgMedsPerLog         MetricChange              `json:"avgMedsPerLog"`
}
