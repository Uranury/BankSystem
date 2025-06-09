package handlers

import (
	"MockBankGo/auth"
	"MockBankGo/middleware"
	"MockBankGo/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	database *sqlx.DB
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{database: db}
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUser models.User
	var existingUser models.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.database.Get(&existingUser, "SELECT * FROM users WHERE email=$1", newUser.Email)
	if err == nil {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	if err != sql.ErrNoRows {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = h.database.Get(&existingUser, "SELECT * FROM users WHERE username=$1", newUser.Username)
	if err == nil {
		http.Error(w, "Username taken", http.StatusConflict)
		return
	}

	if err != sql.ErrNoRows {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	newUser.Password = strings.TrimSpace(newUser.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	newUser.Password = string(hashedPassword)

	_, err = h.database.Exec(
		`INSERT INTO users (username, name, email, password, balance, role) VALUES ($1, $2, $3, $4, $5, $6)`,
		newUser.Username, newUser.Name, newUser.Email, newUser.Password, 0, auth.User,
	)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var logUser models.User
	var existingUser models.User

	if err := json.NewDecoder(r.Body).Decode(&logUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if logUser.Email == "" || logUser.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	err := h.database.Get(&existingUser, "SELECT * FROM users WHERE email = $1", logUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	logUser.Password = strings.TrimSpace(logUser.Password)

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(logUser.Password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(existingUser.ID, string(existingUser.Role))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var users []models.User

	if err := h.database.Select(&users, "SELECT id, username, name, email, balance, role FROM users"); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		http.Error(w, "No users yet", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("JSON encoding error in GetUsers: %v", err)
		return
	}
}

func (h *UserHandler) Profile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, _ := middleware.GetUserID(req.Context())

	var userProfile models.User
	err := h.database.Get(&userProfile, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
