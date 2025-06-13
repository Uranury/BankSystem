package repositories

import (
	"MockBankGo/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ITransationRepository interface {
	WriteTransaction(tx *sqlx.Tx, senderID int64, receiverID int64, amount float64, transactionType models.TransactionType, createdAt time.Time) error
	IncreaseBalance(tx *sqlx.Tx, userID int64, amount float64) error
	DecreaseBalance(tx *sqlx.Tx, userID int64, amount float64) error
	GetTransactions() ([]models.TransactionInfo, error)
}

type TransactionRepository struct {
	database *sqlx.DB
}

func NewTransactionsRepo(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{database: db}
}

func (r *TransactionRepository) WriteTransaction(tx *sqlx.Tx, senderID int64, receiverID int64, amount float64, transactionType models.TransactionType, createdAt time.Time) error {
	_, err := tx.Exec(
		`INSERT INTO transactions (sender_id, receiver_id, amount, type, created_at) VALUES ($1, $2, $3, $4, $5)`,
		senderID, receiverID, amount, transactionType, createdAt,
	)
	return err
}

func (r *TransactionRepository) IncreaseBalance(tx *sqlx.Tx, userID int64, amount float64) error {
	_, err := tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", amount, userID)
	return err
}

func (r *TransactionRepository) DecreaseBalance(tx *sqlx.Tx, userID int64, amount float64) error {
	_, err := tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", amount, userID)
	return err
}

func (r *TransactionRepository) GetTransactions() ([]models.TransactionInfo, error) {
	var transactions []models.TransactionInfo
	err := r.database.Select(&transactions, "SELECT * FROM transactions")
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
