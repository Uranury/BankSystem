package handlers

import (
	"MockBankGo/internal/apperrors"
	"MockBankGo/internal/models"
	"MockBankGo/internal/services"
	"MockBankGo/middleware"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type UserHandler struct {
	database    *sqlx.DB
	userService *services.UserService
}

func NewUserHandler(db *sqlx.DB, user_service *services.UserService) *UserHandler {
	return &UserHandler{database: db, userService: user_service}
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	// handleError now sets its own Content-Type and returns JSON
	w.Header().Set("Content-Type", "application/json")

	if appErr, ok := err.(*apperrors.AppError); ok {
		w.WriteHeader(appErr.StatusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message})
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
	}
}

func (h *UserHandler) Signup(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var NewUser models.User

	if err := json.NewDecoder(req.Body).Decode(&NewUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.userService.CreateUser(&NewUser)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}

func (h *UserHandler) Login(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var logUser models.User

	if err := json.NewDecoder(req.Body).Decode(&logUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.userService.LoginUser(logUser.Email, logUser.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) Profile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := middleware.GetUserID(req.Context())

	userProfile, err := h.userService.GetUserProfile(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := struct {
		Username string  `json:"username"`
		Name     string  `json:"name"`
		Email    string  `json:"email"`
		Balance  float64 `json:"balance"`
	}{
		Username: userProfile.Username,
		Name:     userProfile.Name,
		Email:    userProfile.Email,
		Balance:  userProfile.Balance,
	}

	json.NewEncoder(w).Encode(response)
}
