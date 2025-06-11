package services

import (
	"MockBankGo/internal/apperrors"
	"MockBankGo/internal/models"
	"MockBankGo/internal/repositories"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
	userRepo        *repositories.UserRepository
	database        *sqlx.DB
}

func NewTransactionService(trepo *repositories.TransactionRepository, urepo *repositories.UserRepository, db *sqlx.DB) *TransactionService {
	return &TransactionService{transactionRepo: trepo, userRepo: urepo, database: db}
}

func (s *TransactionService) validateAmount(amount float64) error {
	if amount <= 0 {
		return apperrors.ErrInvalidAmount
	}
	return nil
}

func (s *TransactionService) WithdrawMoney(userID int64, amount float64) error {
	if err := s.validateAmount(amount); err != nil {
		return err
	}

	tx, err := s.database.Beginx()
	if err != nil {
		return apperrors.ErrDatabaseError
	}
	defer tx.Rollback()

	balance, err := s.userRepo.GetBalance(tx, userID)
	if err != nil {
		return apperrors.ErrDatabaseError
	}

	if balance < amount {
		return apperrors.ErrInsufficientFunds
	}

	if err := s.transactionRepo.DecreaseBalance(tx, userID, amount); err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := s.transactionRepo.WriteTransaction(tx, userID, userID, amount, models.Withdraw, time.Now()); err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := tx.Commit(); err != nil {
		return apperrors.ErrDatabaseError
	}

	return nil
}

func (s *TransactionService) DepositMoney(userID int64, amount float64) error {
	if err := s.validateAmount(amount); err != nil {
		return err
	}

	tx, err := s.database.Beginx()
	if err != nil {
		return apperrors.ErrDatabaseError
	}
	defer tx.Rollback()

	if err := s.transactionRepo.IncreaseBalance(tx, userID, amount); err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := s.transactionRepo.WriteTransaction(tx, userID, userID, amount, models.Withdraw, time.Now()); err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := tx.Commit(); err != nil {
		return apperrors.ErrDatabaseError
	}

	return nil
}

func (s *TransactionService) TransferMoney(ctx context.Context, senderID int64, receiverID int64, amount float64) error {
	if err := s.validateAmount(amount); err != nil {
		return err
	}

	if receiverID == senderID {
		return apperrors.ErrSelfTransfer
	}

	tx, err := s.database.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		return apperrors.ErrDatabaseError
	}
	defer tx.Rollback()

	sender, err := s.userRepo.GetUserByIdForUpdate(tx, senderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrSenderNotFound
		}
		return apperrors.ErrDatabaseError
	}

	if sender.Balance < amount {
		return apperrors.ErrInsufficientFunds
	}

	_, err = s.userRepo.GetUserByIdForUpdate(tx, receiverID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrReceiverNotFound
		}
		return apperrors.ErrDatabaseError
	}

	if err := s.transactionRepo.DecreaseBalance(tx, senderID, amount); err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := s.transactionRepo.IncreaseBalance(tx, receiverID, amount); err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := s.transactionRepo.WriteTransaction(tx, senderID, receiverID, amount, models.Transfer, time.Now()); err != nil {
		return apperrors.ErrDatabaseError
	}

	if err := tx.Commit(); err != nil {
		return apperrors.ErrDatabaseError
	}

	return nil
}

func (s *TransactionService) GetTransactions() ([]models.TransactionInfo, error) {
	transactions, err := s.transactionRepo.GetTransactions()
	if err != nil {
		return nil, apperrors.ErrDatabaseError
	}
	return transactions, nil
}
