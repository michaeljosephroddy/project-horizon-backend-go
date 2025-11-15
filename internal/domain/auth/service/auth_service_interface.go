package service

import (
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/model"
)

type IAuthService interface {
	Register(req model.RegisterRequest) (*model.AuthResponse, error)
	Login(req model.LoginRequest) (*model.AuthResponse, error)
	ValidateToken(tokenString string) (uint64, error)
}
