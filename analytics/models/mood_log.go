package models

import "time"

type MoodLog struct {
	MoodLogID  int       `json:"moodLogId"`
	UserID     string    `json:"userId"`
	MoodRating int       `json:"moodRating"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"createdAt"`
	MoodTags   []string  `json:"moodTags"`
}
