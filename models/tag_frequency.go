package models

type TagFrequency struct {
	TagName    string  `json:"tagName"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}
