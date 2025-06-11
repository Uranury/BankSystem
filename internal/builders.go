package internal

import (
	"MockBankGo/internal/handlers"
	"MockBankGo/internal/repositories"
	"MockBankGo/internal/services"

	"github.com/jmoiron/sqlx"
)

func InitUserHandler(db *sqlx.DB) *handlers.UserHandler {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	return handlers.NewUserHandler(db, userService)
}

func InitTransactionHandler(db *sqlx.DB) *handlers.TransactionHandler {
	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionsRepo(db)
	transactionService := services.NewTransactionService(transactionRepo, userRepo, db)
	return handlers.NewTransactionHandler(db, transactionService)
}
