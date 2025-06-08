package handlers

import (
	"MockBankGo/middleware"
	"MockBankGo/models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

type TransactionRequest struct {
	Amount     float64 `json:"amount"`
	ReceiverID *int64  `json:"receiver_id,omitempty"`
}

func (h *UserHandler) Withdraw(w http.ResponseWriter, req *http.Request) {
	var input TransactionRequest
	userId, _ := middleware.GetUserID(req.Context())

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if input.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	tx, _ := h.database.Beginx()
	var balance float64

	err := tx.Get(&balance, "SELECT balance FROM users WHERE id = $1", userId)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Server error", http.StatusInternalServerError)
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
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Withdrawal successful"))
}

func (h *UserHandler) Deposit(w http.ResponseWriter, req *http.Request) {
	var input TransactionRequest
	userID, _ := middleware.GetUserID(req.Context())

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if input.Amount < 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	_, err := h.database.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", input.Amount, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Transfer(w http.ResponseWriter, req *http.Request) {
	var input TransactionRequest
	userID, _ := middleware.GetUserID(req.Context())

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if input.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
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
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var receiver models.User
	var sender models.User

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.Get(&sender, "SELECT * FROM users WHERE id = $1 FOR UPDATE", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Sender not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if sender.Balance < input.Amount {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	err = tx.Get(&receiver, "SELECT * FROM users WHERE id = $1 FOR UPDATE", input.ReceiverID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Receiver not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", input.Amount, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", input.Amount, input.ReceiverID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction successful"))
}
