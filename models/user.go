package models

type User struct {
	ID       int64   `db:"id"`
	Username string  `db:"username"`
	Name     string  `db:"name"`
	Email    string  `db:"email"`
	Password string  `db:"password"`
	Balance  float64 `db:"balance"`
}
