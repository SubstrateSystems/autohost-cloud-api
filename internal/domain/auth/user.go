package auth

import "time"

// User representa un usuario del sistema
type User struct {
	ID           string
	Email        string
	Name         *string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Repository define las operaciones de persistencia para usuarios
type Repository interface {
	CreateUser(email, name, passwordHash string) (string, error)
	FindUserByEmail(email string) (*User, error)
	StoreRefreshToken(userID, tokenHash, userAgent, ip string) error
	RevokeRefreshToken(userID, tokenHash string) error
}
