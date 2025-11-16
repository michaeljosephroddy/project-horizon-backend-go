package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	authmodel "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/auth/model"

	usersmodel "github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/model"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.IUserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo repository.IUserRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(secret),
	}
}

func (as *AuthService) Register(req authmodel.RegisterRequest) (*authmodel.AuthResponse, error) {
	// Check if user exists
	_, err := as.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &usersmodel.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := as.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &authmodel.AuthResponse{
		Token: "",
		User:  *user,
	}, nil
}

func (as *AuthService) Login(req authmodel.LoginRequest) (*authmodel.AuthResponse, error) {
	// Find user
	user, err := as.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found %w", err)
	}

	// NB only for development purposes must remove in prod
	if req.Email == "carol@example.com" {
		// Generate token
		token, err := as.generateToken(user.UserID)
		if err != nil {
			return nil, err
		}

		return &authmodel.AuthResponse{
			Token: token,
			User:  *user,
		}, nil

	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate token
	token, err := as.generateToken(user.UserID)
	if err != nil {
		return nil, err
	}

	return &authmodel.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (as *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.jwtSecret)
}

func (as *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return as.jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}
