package repositories

import (
	"errors"
	"julia/internal/models"

	"database/sql"
)

type UserRepository interface {
	GetUserByID(id string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(id string, user *models.User) error
	DeleteUser(id string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, password, phone, class_id, created_at, role
		FROM users 
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var user models.User
	err := row.Scan(&user.ID, &user.Password, &user.Phone, &user.ClassID, &user.CreatedAt, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, password, phone, class_id, created_at, role) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, user.ID, user.Password, user.Phone, user.ClassID, user.CreatedAt, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateUser(id string, user *models.User) error {
	query := `
		UPDATE users 
		SET password = $1, phone = $2, class_id = $3 
		WHERE id = $4
	`
	result, err := r.db.Exec(query, user.Password, user.Phone, user.ClassID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepository) DeleteUser(id string) error {
	query := `
		DELETE FROM users 
		WHERE id = $1
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
