package models

type User struct {
	ID           string  `db:"id"`
	Email        string  `db:"email"`
	Name         *string `db:"name"`
	PasswordHash string  `db:"password_hash"`
}
