package auth

import "time"

// User representa un usuario del sistema
type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	Name         *string   `db:"name"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Repository define las operaciones de persistencia para usuarios
type Repository interface {
	CreateUser(email, name, passwordHash string) (string, error)
	FindUserByEmail(email string) (*User, error)
	StoreRefreshToken(userID, tokenHash, userAgent, ip string) error
	RevokeRefreshToken(userID, tokenHash string) error
}
