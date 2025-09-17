package contracts

import (
	"github.com/michaeljosephroddy/project-horizon-backend-go/models"
)

type IJournalRepository interface {
	MovingAverages(userId string, startDate string, endDate string) ([]models.MovingAverage)
	MoodTagFrequencies(usersId string, startDate string, endDate string) ([]models.MoodTagFrequency)
	HighDays(userId string, startDate string, endDate string) ([]models.HighDay)
	LowDays(userId string, startDate string, endDate string) ([]models.LowDay)
	StandardDeviation(userId string, startDate string, endDate string) (float32)
}
