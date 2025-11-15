package repository

import (
	"database/sql"
	"errors"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(user *model.User) error {
	query := `
        INSERT INTO user (email, password_hash, created_at, updated_at) 
        VALUES (?, ?, NOW(), NOW())
    `

	result, err := ur.db.Exec(query, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.UserID = int(id)
	return nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	query := `
        SELECT user_id, email, password_hash, created_at, updated_at 
        FROM user 
        WHERE email = ?
    `

	var user model.User
	err := ur.db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return &model.User{}, errors.New("user not found")
	}

	if err != nil {
		return &model.User{}, err
	}

	return &user, nil
}

func (ur *UserRepository) FindByID(id int) (*model.User, error) {
	query := `
        SELECT user_id, email, password_hash, created_at, updated_at 
        FROM user 
        WHERE user_id = ?
    `

	var user model.User
	err := ur.db.QueryRow(query, id).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
