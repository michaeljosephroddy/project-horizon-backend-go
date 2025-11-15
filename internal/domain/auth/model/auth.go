package model

import "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/model"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}
