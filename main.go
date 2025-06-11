package main

import (
	"MockBankGo/db"
	"MockBankGo/internal"
	"MockBankGo/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	router := mux.NewRouter()
	router.Use(middleware.Logger)

	protected := router.NewRoute().Subrouter()
	protected.Use(middleware.JWTAuth)

	userHandler := internal.InitUserHandler(database)
	transactionHandler := internal.InitTransactionHandler(database)

	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/signup", userHandler.Signup).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	protected.HandleFunc("/withdraw", transactionHandler.Withdraw).Methods("POST")
	protected.HandleFunc("/deposit", transactionHandler.Deposit).Methods("POST")
	protected.HandleFunc("/transfer", transactionHandler.Transfer).Methods("POST")
	protected.HandleFunc("/transactions", transactionHandler.GetTransactions).Methods("GET")
	protected.HandleFunc("/profile", userHandler.Profile).Methods("GET")

	log.Printf("Server running on %s", os.Getenv("LISTEN_ADDR"))
	if err := http.ListenAndServe(os.Getenv("LISTEN_ADDR"), router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
