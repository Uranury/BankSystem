package handlers

import (
	"MockBankGo/middleware"
	"MockBankGo/models"
	"encoding/json"
	"net/http"
)

func (h *UserHandler) Withdraw(w http.ResponseWriter, req *http.Request) {
	var input models.TransactionRequest
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
	var input models.TransactionRequest
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
