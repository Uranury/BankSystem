package models

import "MockBankGo/auth"

type User struct {
	ID       int64     `db:"id" json:"id"`             // Changed from int64 to int
	Username string    `db:"username" json:"username"` // Added json tag
	Name     string    `db:"name" json:"name"`         // Added json tag
	Email    string    `db:"email" json:"email"`       // Added json tag
	Password string    `db:"password" json:"password"` // Added json tag
	Balance  float64   `db:"balance" json:"balance"`   // Added json tag
	Role     auth.Role `db:"role" json:"role"`
}
