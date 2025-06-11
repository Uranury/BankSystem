package repositories

import (
	"MockBankGo/internal/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	database *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{database: db}
}

func (r *UserRepository) GetUserById(userId int64) (*models.User, error) {
	var fetchedUser models.User
	err := r.database.Get(&fetchedUser, "SELECT * FROM users WHERE id = $1", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &fetchedUser, nil
}

func (r *UserRepository) GetBalance(tx *sqlx.Tx, userID int64) (float64, error) {
	var balance float64
	err := tx.Get(&balance, "SELECT balance FROM users WHERE id = $1", userID)
	return balance, err
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.database.Select(&users, "SELECT id, username, name, email, balance, role FROM users")
	return users, err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.database.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows // Return nil pointer when no user found
		}
		return nil, err // Return nil pointer for other errors
	}
	return &user, nil // Only return user pointer when actually found
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.database.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	_, err := r.database.Exec(
		`INSERT INTO users (username, name, email, password, balance, role) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.Username, user.Name, user.Email, user.Password, user.Balance, user.Role,
	)
	return err
}

func (r *UserRepository) GetUserByIdForUpdate(tx *sqlx.Tx, userID int64) (*models.User, error) {
	var user models.User
	err := tx.Get(&user, "SELECT * FROM users WHERE id = $1 FOR UPDATE", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}
