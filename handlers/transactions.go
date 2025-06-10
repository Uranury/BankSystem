package handlers

import (
	"MockBankGo/auth"
	"MockBankGo/middleware"
	"MockBankGo/models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionRequest struct {
	Amount     float64 `json:"amount"`
	ReceiverID *int64  `json:"receiver_id,omitempty"`
}

func (h *UserHandler) WriteTransaction(tx *sqlx.Tx, sender_id int64, receiver_id int64, amount float64, transactionType models.TransactionType, created_at time.Time) error {
	_, err := tx.Exec(
		`INSERT INTO transactions (sender_id, receiver_id, amount, type, created_at) VALUES ($1, $2, $3, $4, $5)`,
		sender_id, receiver_id, amount, transactionType, created_at,
	)
	return err
}

func (h *UserHandler) AmountGreaterThanZero(amount float64) bool {
	return amount > 0
}

func (h *UserHandler) Withdraw(w http.ResponseWriter, req *http.Request) {
	var input TransactionRequest
	userId, _ := middleware.GetUserID(req.Context())

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !h.AmountGreaterThanZero(input.Amount) {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	tx, err := h.database.Beginx()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var balance float64

	err = tx.Get(&balance, "SELECT balance FROM users WHERE id = $1", userId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if balance < input.Amount {
		tx.Rollback()
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", input.Amount, userId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.WriteTransaction(tx, userId, userId, input.Amount, models.Withdraw, time.Now())
	if err != nil {
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Withdrawal successful"))
}

func (h *UserHandler) Deposit(w http.ResponseWriter, req *http.Request) {
	var input TransactionRequest
	userID, _ := middleware.GetUserID(req.Context())

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !h.AmountGreaterThanZero(input.Amount) {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	tx, err := h.database.Beginx()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", input.Amount, userID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.WriteTransaction(tx, userID, userID, input.Amount, models.Deposit, time.Now())
	if err != nil {
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deposit successful"))
}

func (h *UserHandler) Transfer(w http.ResponseWriter, req *http.Request) {
	var input TransactionRequest
	userID, _ := middleware.GetUserID(req.Context())

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !h.AmountGreaterThanZero(input.Amount) {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	if input.ReceiverID == nil {
		http.Error(w, "ReceiverID is not provided", http.StatusBadRequest)
		return
	}

	if *input.ReceiverID == userID {
		http.Error(w, "Cannot transfer to yourself", http.StatusBadRequest)
		return
	}

	tx, err := h.database.BeginTxx(req.Context(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var receiver models.User
	var sender models.User

	err = tx.Get(&sender, "SELECT * FROM users WHERE id = $1 FOR UPDATE", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			http.Error(w, "Sender not found", http.StatusNotFound)
			return
		}
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if sender.Balance < input.Amount {
		tx.Rollback()
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	err = tx.Get(&receiver, "SELECT * FROM users WHERE id = $1 FOR UPDATE", input.ReceiverID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			http.Error(w, "Receiver not found", http.StatusNotFound)
			return
		}
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", input.Amount, userID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", input.Amount, input.ReceiverID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.WriteTransaction(tx, userID, *input.ReceiverID, input.Amount, models.Transfer, time.Now())
	if err != nil {
		tx.Rollback()
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction successful"))
}

func (h *UserHandler) GetTransactions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	role, _ := middleware.GetUserRole(req.Context())

	if role != auth.Admin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var transactions []models.TransactionInfo

	err := h.database.Select(&transactions, "SELECT * FROM transactions")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(transactions) == 0 {
		json.NewEncoder(w).Encode([]models.TransactionInfo{})
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
