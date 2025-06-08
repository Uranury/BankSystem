package models

import "time"

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
	Transfer TransactionType = "transfer"
)

type TransactionInfo struct {
	ID         int64           `db:"id" json:"id"`
	SenderID   int64           `db:"sender_id" json:"sender_id"`
	ReceiverID int64           `db:"receiver_id" json:"receiver_id"`
	Amount     float64         `db:"amount" json:"amount"`
	Type       TransactionType `db:"type" json:"type"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
}
