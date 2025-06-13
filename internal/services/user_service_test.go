package services

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"MockBankGo/auth"
	"MockBankGo/internal/apperrors"
	"MockBankGo/internal/models"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	user, ok := args.Get(0).(*models.User)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetUserById(id int64) (*models.User, error) {
	args := m.Called(id)
	user, ok := args.Get(0).(*models.User)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user, ok := args.Get(0).(*models.User)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func TestLoginUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	// Arrange
	password := "securepassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     auth.User,
	}

	// Setup mock: when GetUserByEmail is called, return our user, no error.
	mockRepo.On("GetUserByEmail", user.Email).Return(user, nil)

	service := NewUserService(mockRepo)

	// Act
	token, err := service.LoginUser(user.Email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)

	password := "securepassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     auth.User,
	}

	mockRepo.On("GetUserByEmail", user.Email).Return(user, nil)

	service := NewUserService(mockRepo)

	// Try wrong password
	token, err := service.LoginUser(user.Email, "wrongpassword")

	assert.Equal(t, apperrors.ErrInvalidCredentials, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	email := "nonexistent@example.com"

	// Return nil user and an sql.ErrNoRows-like behavior: just nil user and no error triggers user not found logic.
	mockRepo.On("GetUserByEmail", email).Return(nil, sql.ErrNoRows)

	service := NewUserService(mockRepo)

	token, err := service.LoginUser(email, "anyPassword")

	assert.Equal(t, apperrors.ErrUserNotFound, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginUser_EmptyInput(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	token, err := service.LoginUser("", "")

	assert.Equal(t, apperrors.ErrEmailPasswordRequired, err)
	assert.Empty(t, token)

	// No repo method should be called in this case
	mockRepo.AssertExpectations(t)
}
