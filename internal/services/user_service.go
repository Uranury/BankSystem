package services

import (
	"MockBankGo/auth"
	"MockBankGo/internal/apperrors"
	"MockBankGo/internal/models"
	"MockBankGo/internal/repositories"
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.IUserRepository
}

func NewUserService(UserRepo repositories.IUserRepository) *UserService {
	return &UserService{userRepo: UserRepo}
}

func (s *UserService) CreateUser(user *models.User) error {
	existingUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil && err != sql.ErrNoRows {
		return apperrors.ErrDatabaseError
	}
	if existingUser != nil {
		return apperrors.ErrEmailAlreadyInUse
	}

	existingUser, err = s.userRepo.GetUserByUsername(user.Username)
	if err != nil && err != sql.ErrNoRows {
		return apperrors.ErrDatabaseError
	}
	if existingUser != nil {
		return apperrors.ErrUsernameAlreadyInUse
	}

	user.Password = strings.TrimSpace(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.ErrHashPassword
	}

	user.Password = string(hashedPassword)
	user.Balance = 0
	user.Role = auth.User

	if err := s.userRepo.CreateUser(user); err != nil {
		return apperrors.ErrDatabaseError
	}

	return nil
}

func (s *UserService) LoginUser(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", apperrors.ErrEmailPasswordRequired
	}

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apperrors.ErrUserNotFound
		}
		return "", apperrors.ErrDatabaseError
	}

	password = strings.TrimSpace(password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := auth.GenerateJWT(user.ID, string(user.Role))
	if err != nil {
		return "", apperrors.ErrGenerateToken
	}

	return token, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.userRepo.GetUsers()
	if err != nil {
		return nil, apperrors.ErrDatabaseError
	}

	if len(users) == 0 {
		return nil, apperrors.ErrNoUsersFound
	}

	return users, nil
}

func (s *UserService) GetUserProfile(userID int64) (*models.User, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}

	user.Password = ""
	return user, nil
}
