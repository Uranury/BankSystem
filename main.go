package main

import (
	"MockBankGo/db"
	"MockBankGo/handlers"
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
	protected := router.NewRoute().Subrouter()
	protected.Use(middleware.JWTAuth)

	userHandler := handlers.NewUserHandler(database)
	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/signup", userHandler.Signup).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	protected.HandleFunc("/withdraw", userHandler.Withdraw).Methods("POST")
	protected.HandleFunc("/deposit", userHandler.Deposit).Methods("POST")

	log.Printf("Server running on %s", os.Getenv("LISTEN_ADDR"))
	if err := http.ListenAndServe(os.Getenv("LISTEN_ADDR"), router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
