package models

import "time"

type MovingAverage struct {
	Date      time.Time `json:"date"`
	MovingAvg float64   `json:"movingAvg"`
}
