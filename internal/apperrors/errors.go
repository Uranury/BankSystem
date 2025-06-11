package apperrors

import "net/http"

type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(message string, statuscode int) *AppError {
	return &AppError{Message: message, StatusCode: statuscode}
}

// Predefined errors for User operations
var (
	ErrEmailAlreadyInUse = &AppError{
		Message:    "email already in use",
		StatusCode: http.StatusConflict,
	}

	ErrUsernameAlreadyInUse = &AppError{
		Message:    "username already in use",
		StatusCode: http.StatusConflict,
	}

	ErrUserNotFound = &AppError{
		Message:    "user not found",
		StatusCode: http.StatusNotFound,
	}

	ErrInvalidCredentials = &AppError{
		Message:    "incorrect password",
		StatusCode: http.StatusUnauthorized,
	}

	ErrEmailPasswordRequired = &AppError{
		Message:    "email and password are required",
		StatusCode: http.StatusBadRequest,
	}

	ErrNoUsersFound = &AppError{
		Message:    "no users yet",
		StatusCode: http.StatusNotFound,
	}
)

// Predefined errors for Transaction operations
var (
	ErrInvalidAmount = &AppError{
		Message:    "amount must be greater than zero",
		StatusCode: http.StatusBadRequest,
	}

	ErrInsufficientFunds = &AppError{
		Message:    "insufficient funds",
		StatusCode: http.StatusBadRequest,
	}

	ErrSelfTransfer = &AppError{
		Message:    "cannot transfer to yourself",
		StatusCode: http.StatusBadRequest,
	}

	ErrReceiverNotFound = &AppError{
		Message:    "receiver not found",
		StatusCode: http.StatusNotFound,
	}

	ErrSenderNotFound = &AppError{
		Message:    "sender not found",
		StatusCode: http.StatusNotFound,
	}

	ErrReceiverNotProvided = &AppError{
		Message:    "receiver not provided",
		StatusCode: http.StatusBadRequest,
	}
)

// Generic server errors
var (
	ErrInternalServer = &AppError{
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	}

	ErrDatabaseError = &AppError{
		Message:    "database error",
		StatusCode: http.StatusInternalServerError,
	}

	ErrHashPassword = &AppError{
		Message:    "failed to hash password",
		StatusCode: http.StatusInternalServerError,
	}

	ErrGenerateToken = &AppError{
		Message:    "failed to generate token",
		StatusCode: http.StatusInternalServerError,
	}
)
