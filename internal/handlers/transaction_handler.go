package handlers

import (
	"MockBankGo/internal/apperrors"
	"MockBankGo/internal/services"
	"MockBankGo/middleware"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type TransactionHandler struct {
	database           *sqlx.DB
	transactionService *services.TransactionService
}

type TransactionInput struct {
	Amount     float64 `json:"amount"`
	ReceiverID *int64  `json:"receiver_id,omitempty"`
}

func NewTransactionHandler(db *sqlx.DB, transaction_service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{database: db, transactionService: transaction_service}
}

func (h *TransactionHandler) handleError(w http.ResponseWriter, err error) {
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

func (h *TransactionHandler) Withdraw(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input TransactionInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID, _ := middleware.GetUserID(req.Context())

	err := h.transactionService.WithdrawMoney(userID, input.Amount)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Write([]byte("Withdrawal successful!"))
}

func (h *TransactionHandler) Deposit(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input TransactionInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID, _ := middleware.GetUserID(req.Context())

	err := h.transactionService.DepositMoney(userID, input.Amount)
	if err != nil {
		h.handleError(w, err)
		return
	}
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var input TransactionInput

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.ReceiverID == nil {
		h.handleError(w, apperrors.ErrReceiverNotProvided)
		return
	}

	userID, _ := middleware.GetUserID(req.Context())

	if err := h.transactionService.TransferMoney(req.Context(), userID, *input.ReceiverID, input.Amount); err != nil {
		h.handleError(w, err)
		return
	}

	w.Write([]byte("Transfer successful!"))
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transactions, err := h.transactionService.GetTransactions()
	if err != nil {
		h.handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(transactions)
}
