package repository

import (
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/model"
)

type IUserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id int) (*model.User, error)
}
