package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/michaeljosephroddy/project-horizon-backend-go/internal/domain/users/model"
	"os"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

var home = os.Getenv("HOME")
var queriesDir = "/repos/project-horizon-backend-go/internal/domain/users/repository/queries/"

func (ur *UserRepository) Create(user *model.User) error {
	fileName := "create.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

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
	fileName := "find_by_email.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return &model.User{}, fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	var user model.User
	var createdAtStr, updatedAtStr string
	queryErr := ur.db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&createdAtStr,
		&updatedAtStr,
	)
	if queryErr == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if queryErr != nil {
		return nil, err
	}
	// Parse MySQL DATETIME format
	layout := "2006-01-02 15:04:05"
	user.CreatedAt, err = time.Parse(layout, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed parsing created_at: %w", err)
	}
	user.UpdatedAt, err = time.Parse(layout, updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed parsing updated_at: %w", err)
	}
	return &user, nil
}

func (ur *UserRepository) FindByID(id int) (*model.User, error) {
	fileName := "find_by_id.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return &model.User{}, fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	var user model.User
	queryErr := ur.db.QueryRow(query, id).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if queryErr == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if queryErr != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserMedicationsByUserID retrieves all active medications for a user
func (ur *UserRepository) GetMedicationsByUserID(userID int) ([]model.UserMedicationDTO, error) {

	fileName := "get_medications_by_user_id.sql"
	path := home + queriesDir + fileName

	content, err := os.ReadFile(path)
	if err != nil {
		return make([]model.UserMedicationDTO, 0), fmt.Errorf("error reading query file %w", err)
	}
	query := string(content)

	rows, err := ur.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user medications: %w", err)
	}
	defer rows.Close()

	var medications []model.UserMedicationDTO
	for rows.Next() {
		var med model.UserMedicationDTO
		var startDate string

		err := rows.Scan(
			&med.UserMedicationID,
			&med.MedicationID,
			&med.Dosage,
			&startDate,
			&med.Note,
			&med.Name,
			&med.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan medication row: %w", err)
		}

		// Parse the date string to time.Time
		parsedDate, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse moving average date: %w", err)
		}

		// Format date as string for JSON response
		med.StartDate = parsedDate
		medications = append(medications, med)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating medication rows: %w", err)
	}

	// Return empty slice instead of nil if no medications found
	if medications == nil {
		medications = []model.UserMedicationDTO{}
	}

	return medications, nil
}
