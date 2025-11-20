// internal/domain/users/service/users_service.go
package service

import (
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/repository"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/model"
)

type UserService struct {
	usersRepo repository.IUserRepository
}

func NewUsersService(usersRepo repository.IUserRepository) *UserService {
	return &UserService{
		usersRepo: usersRepo,
	}
}

func (us *UserService) GetUserMedications(userID int) ([]model.UserMedicationDTO, error) {
	return us.usersRepo.GetMedicationsByUserID(userID)
}
