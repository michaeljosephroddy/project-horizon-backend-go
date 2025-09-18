package models

type MovingAverage struct {
	Date string `json:"date"`
	ThreeDay float32 `json:"threeDayMovingAvg"` 
}
